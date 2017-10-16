package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.Cli;
import com.marshmallow.anwork.app.cli.Command;
import com.marshmallow.anwork.app.cli.Flag;
import com.marshmallow.anwork.app.cli.List;
import com.marshmallow.anwork.app.cli.Visitor;

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
public class AnworkCliDocumentationGenerator implements Visitor {

  private static final String NO_DESCRIPTION = "<no description>";
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
      Visitor visitor = new AnworkCliDocumentationGenerator(writer);
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

  @Override
  public void visitFlag(Flag flag) {
    checkFlagState();
    StringBuilder lineBuilder = new StringBuilder();
    lineBuilder.append("- -").append(flag.getShortFlag());
    if (flag.hasLongFlag()) {
      lineBuilder.append("|--" + flag.getLongFlag());
    }
    if (flag.hasArgument()) {
      lineBuilder.append(" (")
                 .append(flag.getArgument().getType().toString())
                 .append(" ")
                 .append(flag.getArgument().getName());
      if (flag.getArgument().hasDescription()) {
        lineBuilder.append(": ")
                   .append(flag.getArgument().getDescription());
      }
      lineBuilder.append(")");
    }
    lineBuilder.append(": ").append((flag.hasDescription()
                                     ? flag.getDescription()
                                     : NO_DESCRIPTION));
    writer.println(lineBuilder.toString());
  }

  @Override
  public void visitList(List list) {
    checkListState();
    String prefix = listStack.stream().collect(Collectors.joining(" "));
    writer.println(String.format("# %s *%s*: %s",
                                 prefix,
                                 list.getName(),
                                 (list.hasDescription()
                                  ? list.getDescription()
                                  : NO_DESCRIPTION)));
    listStack.push(list.getName());
  }

  @Override
  public void leaveList(List list) {
    checkListState();
    listStack.pop();
  }

  @Override
  public void visitCommand(Command command) {
    checkCommandState();
    String prefix = listStack.stream().collect(Collectors.joining(" "));
    writer.println(String.format("- %s *%s*: %s",
                                 prefix,
                                 command.getName(),
                                 (command.hasDescription()
                                  ? command.getDescription()
                                  : NO_DESCRIPTION)));
  }
}
