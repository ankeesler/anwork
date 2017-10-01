package com.marshmallow.anwork.app.cli;

import org.xml.sax.helpers.DefaultHandler;

import java.util.Stack;

import org.xml.sax.Attributes;
import org.xml.sax.SAXException;
import org.xml.sax.SAXParseException;

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
  private static final String FLAG_PARAMETERNAME = "parameterName";
  private static final String FLAG_DESCRIPTION = "description";
  private static final String FLAG_ACTIONCREATOR = "actionCreator";

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
  public void startElement(String uri, String localName, String qName, Attributes attributes) {
    debug("startElement(uri=" + uri
                     + ", localName=" + localName
                     + ", qName=" + qName
                     + ", attributes=" + attributes + ")");
    if (qName.equals(CLI)) {
      String name = attributes.getValue(CLI_NAME);
      String description = attributes.getValue(CLI_DESCRIPTION);
      makeCli(name, description);
    } else if (qName.equals(FLAG)) {
      String shortFlag = attributes.getValue(FLAG_SHORTFLAG);
      String longFlag = attributes.getValue(FLAG_LONGFLAG);
      String parameterName = attributes.getValue(FLAG_PARAMETERNAME);
      String description = attributes.getValue(FLAG_DESCRIPTION);
      String actionCreator = attributes.getValue(FLAG_ACTIONCREATOR);
      makeFlag(shortFlag, longFlag, parameterName, description, actionCreator);
    } else if (qName.equals(COMMAND)) {
      String name = attributes.getValue(COMMAND_NAME);
      String description = attributes.getValue(COMMAND_DESCRIPTION);
      String actionCreator = attributes.getValue(COMMAND_ACTIONCREATOR);
      makeCommand(name, description, actionCreator);
    } else if (qName.equals(LIST)) {
      String name = attributes.getValue(LIST_NAME);
      String description = attributes.getValue(LIST_DESCRIPTION);
      makeList(name, description);
    }
  }

  @Override
  public void endElement(String uri, String localName, String qName) {
    debug("endElement(uri=" + uri
                   + ", localName=" + localName
                   + ", qName=" + qName + ")");
    if (qName.equals(LIST)) {
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

  private void makeCli(String name, String description) {
    cli = new Cli(name, description);
    listStack.push(cli.getRoot());
  }

  private void makeFlag(String shortFlag,
                        String longFlag,
                        String parameterName,
                        String description,
                        String actionCreator) {
    CliAction realAction = new CliUsageAction((CliNodeImpl)cli.getRoot()); // TODO!
    if (longFlag == null) {
      if (parameterName == null) {
        listStack.peek().addShortFlag(shortFlag, description, realAction);
      } else {
        listStack.peek().addShortFlagWithParameter(shortFlag, description, parameterName, realAction);
      }
    } else {
      if (parameterName == null) {
        listStack.peek().addLongFlag(shortFlag, longFlag, description, realAction);
      } else {
        listStack.peek().addLongFlagWithParameter(shortFlag, longFlag, description, parameterName, realAction);
      }
    }
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
