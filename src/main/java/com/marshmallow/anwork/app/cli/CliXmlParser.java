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

  // <flag>
  private static final String FLAG = "flag";
  private static final String FLAG_SHORTFLAG = "shortFlag";
  private static final String FLAG_LONGFLAG = "longFlag";
  private static final String FLAG_DESCRIPTION = "description";
  private static final String FLAG_PARAMETER = "parameter";
  private static final String FLAG_PARAMETER_NAME = "name";
  private static final String FLAG_PARAMETER_DESCRIPTION = "description";
  private static final String FLAG_PARAMETER_TYPE = "type";

  // <command>
  private static final String COMMAND = "command";
  private static final String COMMAND_NAME = "name";
  private static final String COMMAND_DESCRIPTION = "description";
  private static final String COMMAND_ACTIONCREATOR = "actionCreator";

  // <list>
  private static final String LIST = "list";
  private static final String LIST_NAME = "name";
  private static final String LIST_DESCRIPTION = "description";

  // <cli>
  private static final String CLI = "cli";
  private static final String CLI_NAME = LIST_NAME;
  private static final String CLI_DESCRIPTION = LIST_DESCRIPTION;

  private static final boolean DEBUG = true;

  private static void debug(String string) {
    if (DEBUG) {
      System.out.println(string);
    }
  }

  private Cli cli;
  private Stack<CliList> listStack = new Stack<CliList>();

  // TODO: make this less field-driven and use a builder paradigm!
  private String flagShortFlagName;
  private String flagLongFlagName;
  private String flagDescription;
  private String flagParameterName;
  private String flagParameterDescription;
  private String flagParameterType;

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
                           Attributes attributes) {
    debug("startElement(uri=" + uri
                     + ", localName=" + localName
                     + ", elementName=" + elementName
                     + ", attributes=" + attributes + ")");
    if (elementName.equals(CLI)) {
      String name = attributes.getValue(CLI_NAME);
      String description = attributes.getValue(CLI_DESCRIPTION);
      makeCli(name, description);
    } else if (elementName.equals(FLAG)) {
      flagShortFlagName = attributes.getValue(FLAG_SHORTFLAG);
      flagLongFlagName = attributes.getValue(FLAG_LONGFLAG);
      flagDescription = attributes.getValue(FLAG_DESCRIPTION);
    } else if (elementName.equals(FLAG_PARAMETER)) {
      flagParameterName = attributes.getValue(FLAG_PARAMETER_NAME);
      flagParameterDescription = attributes.getValue(FLAG_PARAMETER_DESCRIPTION);
      flagParameterType = attributes.getValue(FLAG_PARAMETER_TYPE);
    } else if (elementName.equals(COMMAND)) {
      String name = attributes.getValue(COMMAND_NAME);
      String description = attributes.getValue(COMMAND_DESCRIPTION);
      String actionCreator = attributes.getValue(COMMAND_ACTIONCREATOR);
      makeCommand(name, description, actionCreator);
    } else if (elementName.equals(LIST)) {
      String name = attributes.getValue(LIST_NAME);
      String description = attributes.getValue(LIST_DESCRIPTION);
      makeList(name, description);
    }
  }

  @Override
  public void endElement(String uri, String localName, String elementName) {
    debug("endElement(uri=" + uri
                   + ", localName=" + localName
                   + ", elementName=" + elementName + ")");
    if (elementName.equals(LIST)) {
      listStack.pop();
    } else if (elementName.equals(FLAG)) {
      makeFlag();
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

  private void makeCli(String name, String description) {
    cli = new Cli(name, description);
    listStack.push(cli.getRoot());
  }

  private void makeFlag() {
    // See note in cli.xsd - it should be mandated by the schema that these values stay up to date
    // with the CliArgumentType enum!
    boolean hasParameter = (flagParameterName != null);
    CliArgumentType parameterType = (hasParameter
                                     ? CliArgumentType.valueOf(flagParameterType)
                                     : null);
    if (flagLongFlagName == null) {
      if (!hasParameter) {
        listStack.peek().addShortFlag(flagShortFlagName, flagDescription);
      } else {
        listStack.peek().addShortFlagWithParameter(flagShortFlagName,
                                                   flagDescription,
                                                   flagParameterName,
                                                   flagParameterDescription,
                                                   parameterType);
      }
    } else {
      if (!hasParameter) {
        listStack.peek().addLongFlag(flagShortFlagName, flagLongFlagName, flagDescription);
      } else {
        listStack.peek().addLongFlagWithParameter(flagShortFlagName,
                                                  flagLongFlagName,
                                                  flagDescription,
                                                  flagParameterName,
                                                  flagParameterDescription,
                                                  parameterType);
      }
    }
    flagShortFlagName = null;
    flagLongFlagName = null;
    flagDescription = null;
    flagParameterName = null;
    flagDescription = null;
    flagParameterType = null;
  }

  private void makeCommand(String name,
                           String description,
                           String action) {
    CliAction realAction = new CliUsageAction((CliNodeImpl)cli.getRoot()); // TODO!
    listStack.peek().addCommand(name, description, realAction);
  }

  private void makeList(String name,
                        String description) {
    listStack.push(listStack.peek().addList(name, description));
  }
}
