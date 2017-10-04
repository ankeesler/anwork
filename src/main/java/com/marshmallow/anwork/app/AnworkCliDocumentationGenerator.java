package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.Cli;
import com.marshmallow.anwork.app.cli.CliArgumentType;
import com.marshmallow.anwork.app.cli.CliVisitor;

import java.io.PrintWriter;
import java.util.Stack;
import java.util.stream.Collectors;

/**
 * This is an app that generates the documentation for the ANWORK CLI.
 *
 * <p>
 * Created Oct 1, 2017
 * </p>
 *
 * @author Andrew
 */
public class AnworkCliDocumentationGenerator implements CliVisitor {

  private static final String FILENAME = "doc/CLI.md";

  private static enum State {
    NONE,
    FLAG,
    LIST,
    COMMAND,
    ;
  }

  /**
   * This is the main method for the {@link AnworkCliDocumentationGenerator} app.
   *
   * @param args The command line arguments for this app.
   */
  public static void main(String[] args) {
    try (PrintWriter writer = new PrintWriter(FILENAME)) {
      Cli cli = AnworkApp.createCli();
      CliVisitor visitor = new AnworkCliDocumentationGenerator(writer);
      cli.visit(visitor);
    } catch (Exception e) {
      System.out.println("Error: " + e);
    }
  }

  private PrintWriter writer;
  private State state;
  private Stack<String> listStack;

  private AnworkCliDocumentationGenerator(PrintWriter writer) {
    this.writer = writer;
    this.state = State.NONE;
    this.listStack = new Stack<String>();

    writer.println("This documentation is generated from " + this.getClass().getName());
    writer.println();
  }

  private void checkFlagState() {
    if (state != State.FLAG) {
      writer.println("## Flags");
      state = State.FLAG;
    }
  }

  private void checkListState() {
    if (state != State.LIST) {
      state = State.LIST;
    }
  }

  private void checkCommandState() {
    if (state != State.COMMAND) {
      writer.println("## Commands");
      state = State.COMMAND;
    }
  }

  private void writeFlagLine(String shortFlag,
                             String longFlag,
                             String description,
                             String parameterName,
                             String parameterDescription,
                             CliArgumentType parameterType) {
    StringBuilder lineBuilder = new StringBuilder();
    lineBuilder.append("- -").append(shortFlag);
    if (longFlag != null) {
      lineBuilder.append("|--" + longFlag);
    }
    if (parameterName != null) {
      lineBuilder.append(" (")
                 .append(parameterType.name())
                 .append(" ")
                 .append(parameterName)
                 .append(": ")
                 .append(parameterDescription)
                 .append(")");
    }
    lineBuilder.append(": ").append(description);
    writer.println(lineBuilder.toString());
  }

  @Override
  public void visitShortFlag(String shortFlag, String description) {
    checkFlagState();
    writeFlagLine(shortFlag, null, description, null, null, null);
  }

  @Override
  public void visitShortFlagWithParameter(String shortFlag,
                                          String description,
                                          String parameterName,
                                          String parameterDescription,
                                          CliArgumentType parameterType) {
    checkFlagState();
    writeFlagLine(shortFlag,
                  null,
                  description,
                  parameterName,
                  parameterDescription,
                  parameterType);
  }

  @Override
  public void visitLongFlag(String shortFlag,
                            String longFlag,
                            String description) {
    checkFlagState();
    writeFlagLine(shortFlag, longFlag, description, null, null, null);
  }

  @Override
  public void visitLongFlagWithParameter(String shortFlag,
                                         String longFlag,
                                         String description,
                                         String parameterName,
                                         String parameterDescription,
                                         CliArgumentType parameterType) {
    checkFlagState();
    writeFlagLine(shortFlag,
                  longFlag,
                  description,
                  parameterName,
                  parameterDescription,
                  parameterType);
  }

  @Override
  public void visitList(String name, String description) {
    checkListState();
    String prefix = listStack.stream().collect(Collectors.joining(" "));
    writer.println(String.format("# %s *%s*: %s", prefix, name, description));
    listStack.push(name);
  }

  @Override
  public void leaveList(String name) {
    checkListState();
    listStack.pop();
  }

  @Override
  public void visitCommand(String name, String description) {
    checkCommandState();
    String prefix = listStack.stream().collect(Collectors.joining(" "));
    writer.println(String.format("- %s *%s*: %s", prefix, name, description));
  }
}
