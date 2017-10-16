package com.marshmallow.anwork.app.cli;

import java.util.Stack;

import org.xml.sax.Attributes;
import org.xml.sax.SAXException;
import org.xml.sax.SAXParseException;
import org.xml.sax.helpers.DefaultHandler;

/**
 * This class is a {@link DefaultHandler} for use by the {@link CliXmlReader}.
 *
 * <p>
 * Created Oct 1, 2017
 * </p>
 *
 * @author Andrew
 */
class CliXmlParser extends DefaultHandler {

  // <argument>
  private static final String ARGUMENT = "argument";
  private static final String ARGUMENT_NAME = "name";
  private static final String ARGUMENT_DESCRIPTION = "description";
  private static final String ARGUMENT_TYPE = "type";

  // <flag>
  private static final String FLAG = "flag";
  private static final String FLAG_SHORTFLAG = "shortFlag";
  private static final String FLAG_LONGFLAG = "longFlag";
  private static final String FLAG_DESCRIPTION = "description";

  // <command>
  private static final String COMMAND = "command";
  private static final String COMMAND_NAME = "name";
  private static final String COMMAND_DESCRIPTION = "description";

  // <action>
  private static final String ACTION = "action";
  private static final String ACTION_CLASS = "class";

  // <actionCreator>
  private static final String ACTIONCREATOR = "actionCreator";
  private static final String ACTIONCREATOR_CLASS = "class";

  // <list>
  private static final String LIST = "list";
  private static final String LIST_NAME = "name";
  private static final String LIST_DESCRIPTION = "description";

  // <cli>
  private static final String CLI = "cli";
  private static final String CLI_NAME = LIST_NAME;
  private static final String CLI_DESCRIPTION = LIST_DESCRIPTION;

  private static final boolean DEBUG = false;

  private static void debug(String string) {
    if (DEBUG) {
      System.out.println(string);
    }
  }

  private Cli cli;
  private Stack<MutableList> listStack = new Stack<MutableList>();
  private MutableFlag currentFlag;
  private MutableCommand currentCommand;

  /**
   * Get the parsed {@link Cli} object. This method is only valid once the parsing has taken place!
   *
   * @return The parsed {@link Cli} object
   */
  public Cli getCli() {
    return cli;
  }

  @Override
  public void startDocument() {
    debug("startDocument");
  }

  @Override
  public void endDocument() {
    debug("endDocument");
  }

  @Override
  public void startElement(String uri,
                           String localName,
                           String elementName,
                           Attributes attributes) throws SAXException {
    debug("startElement(uri=" + uri
                     + ", localName=" + localName
                     + ", elementName=" + elementName
                     + ", attributes=" + attributes + ")");
    if (elementName.equals(CLI)) {
      String name = attributes.getValue(CLI_NAME);
      String description = attributes.getValue(CLI_DESCRIPTION);
      cli = makeCli(name, description);
      listStack.push(cli.getRoot());
    } else if (elementName.equals(ARGUMENT)) {
      String name = attributes.getValue(ARGUMENT_NAME);
      String description = attributes.getValue(ARGUMENT_DESCRIPTION);
      String type = attributes.getValue(ARGUMENT_TYPE);
      addArgument(currentFlag, name, description, type);
    } else if (elementName.equals(FLAG)) {
      String shortFlag = attributes.getValue(FLAG_SHORTFLAG);
      String longFlag = attributes.getValue(FLAG_LONGFLAG);
      String description = attributes.getValue(FLAG_DESCRIPTION);
      currentFlag = addFlag(currentCommand != null ? currentCommand : listStack.peek(),
                            shortFlag,
                            longFlag,
                            description);
    } else if (elementName.equals(COMMAND)) {
      String name = attributes.getValue(COMMAND_NAME);
      String description = attributes.getValue(COMMAND_DESCRIPTION);
      currentCommand = addCommand(listStack.peek(), name, description);
    } else if (elementName.equals(ACTION) || elementName.equals(ACTIONCREATOR)) {
      boolean isActionCreator = elementName.equals(ACTIONCREATOR);
      String className = attributes.getValue(isActionCreator ? ACTIONCREATOR_CLASS : ACTION_CLASS);
      try {
        setAction(currentCommand, className, isActionCreator);
      } catch (Exception e) {
        throw new SAXException(e);
      }
    } else if (elementName.equals(LIST)) {
      String name = attributes.getValue(LIST_NAME);
      String description = attributes.getValue(LIST_DESCRIPTION);
      listStack.push(addList(listStack.peek(), name, description));
    }
  }

  @Override
  public void endElement(String uri, String localName, String elementName) throws SAXException {
    debug("endElement(uri=" + uri
                   + ", localName=" + localName
                   + ", elementName=" + elementName + ")");
    if (elementName.equals(FLAG)) {
      currentFlag = null;
    } else if (elementName.equals(COMMAND)) {
      currentCommand = null;
    } else if (elementName.equals(LIST)) {
      listStack.pop();
    }
  }

  @Override
  public void warning(SAXParseException e) {
    debug("warning(" + e + ")");
  }

  @Override
  public void error(SAXParseException e) throws SAXException {
    debug("error(" + e + ")");
    throw e;
  }

  @Override
  public void fatalError(SAXParseException e) throws SAXException {
    debug("fatalError(" + e + ")");
    throw e;
  }

  private static Cli makeCli(String name, String description) {
    Cli cli = new Cli(name);
    if (description != null) {
      cli.getRoot().setDescription(description);
    }
    return cli;
  }

  private static ArgumentType<?> getArgumentType(String type) {
    // TODO: get rid of this horrible hardcoding!
    switch (type) {
      case "STRING":
        return ArgumentType.STRING;
      case "NUMBER":
        return ArgumentType.NUMBER;
      case "BOOLEAN":
        return ArgumentType.BOOLEAN;
      default:
        throw new IllegalArgumentException("No known ArgumentType for type" + type);
    }
  }

  private static void addArgument(MutableFlag flag, String name, String description, String type) {
    ArgumentType<?> realType = getArgumentType(type);
    MutableArgument argument = flag.setArgument(name, realType);
    if (description != null) {
      argument.setDescription(description);
    }
  }

  private static MutableFlag addFlag(MutableListOrCommand listOrCommand,
                                     String shortFlag,
                                     String longFlag,
                                     String description) {
    MutableFlag flag = listOrCommand.addFlag(shortFlag);
    if (longFlag != null) {
      flag.setLongFlag(longFlag);
    }
    if (description != null) {
      flag.setDescription(description);
    }
    return flag;
  }

  private static MutableCommand addCommand(MutableList list, String name, String description) {
    // Action is set later...
    MutableCommand command = list.addCommand(name, null);
    if (description != null) {
      command.setDescription(description);
    }
    return command;
  }

  private static void setAction(MutableCommand command,
                                String className,
                                boolean isActionCreator) throws Exception {
    String commandName = command.getName();
    Class<?> clazz = Class.forName(className);
    if ((isActionCreator && !ActionCreator.class.isAssignableFrom(clazz))
        || (!isActionCreator && !Action.class.isAssignableFrom(clazz))) {
      throw new Exception("Class " + clazz.getName() + " for command " + commandName
                          + " is not an instance of " + Action.class.getName());
    }
    Action action =  (isActionCreator
                      ? ((ActionCreator)clazz.newInstance()).createAction(commandName)
                      : (Action)clazz.newInstance());
    command.setAction(action);
  }

  private static MutableList addList(MutableList currentList,
                                     String name,
                                     String description) {
    MutableList list = currentList.addList(name);
    if (description != null) {
      list.setDescription(description);
    }
    return list;
  }
}
