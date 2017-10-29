package com.marshmallow.anwork.app.cli;

import java.io.PrintWriter;
import java.util.Stack;
import java.util.stream.Collectors;

/**
 * This is a {@link DocumentationGenerator} that creates {@link Cli} API documentation in
 * Github README format.
 *
 * <p>
 * Created Oct 16, 2017
 * </p>
 *
 * @author Andrew
 */
class GithubReadmeDocumentationGenerator implements DocumentationGenerator, Visitor {

  private static final String NO_DESCRIPTION = "<no description>";

  private static enum State {
    NONE,
    FLAG,
    LIST,
    COMMAND,
    ;
  }

  private PrintWriter writer;
  private State state;
  private Stack<String> listStack;

  @Override
  public DocumentationType getType() {
    return DocumentationType.GITHUB_MARKDOWN;
  }

  @Override
  public void generate(Cli cli, PrintWriter writer) {
    this.writer = writer;
    this.state = State.NONE;
    this.listStack = new Stack<String>();

    writer.println("This documentation is generated from " + getClass().getName());
    writer.println();
    cli.visit(this);
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
                 .append(flag.getArgument().getType().getName())
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
