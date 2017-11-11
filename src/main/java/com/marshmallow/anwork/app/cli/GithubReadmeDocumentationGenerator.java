package com.marshmallow.anwork.app.cli;

import java.io.PrintWriter;
import java.util.Stack;
import java.util.stream.Collectors;

/**
 * This is a {@link DocumentationGenerator} that creates {@link Cli} API documentation in
 * Github README format.
 *
 * <p>
 * The format of this documentation is roughly the following.
 * <pre>
 * `root-list [-d|--debug] ...` : root-list description
 * {@literal *} `root-list command-a {@literal <}foo:STRING{@literal >} {@literal <}bar:NUMB...`
 *
 * `root-list sub-list-a ...`: sub-list-a description
 * {@literal *} `root-list sub-list-a command-a {@literal <}foo:STRING{@literal >}`
 * {@literal *} {@literal *} command-a description
 * {@literal *} `root-list sub-list-a command-b [-a--andrew {@literal <}foo:STRING{@literal >}]`
 * {@literal *} {@literal *} command-b description
 * {@literal *} `root-list sub-list-a command-c {@literal <}foo:STRING{@literal >} {@literal <}...`
 * {@literal *} {@literal *} command-c description
 *
 * `root-list sub-list-b [-C {@literal <}directory:STRING{@literal >}] ...`: sub-list-b description
 * {@literal *} `root-list sub-list-b command-a {@literal <}foo:STRING{@literal >}`
 * {@literal *} {@literal *} command-a description
 * {@literal *} `root-list sub-list-b command-b`
 * {@literal *} {@literal *} command-b description
 * {@literal *} `root-list sub-list-b command-c {@literal <}foo:STRING{@literal >} {@literal <}...`
 * {@literal *} {@literal *} command-c description
 * </pre>
 * </p>
 *
 * <p>
 * Created Oct 16, 2017
 * </p>
 *
 * @author Andrew
 */
class GithubReadmeDocumentationGenerator implements DocumentationGenerator, Visitor {

  private PrintWriter writer;
  private Stack<String> listStack;

  @Override
  public DocumentationType getType() {
    return DocumentationType.GITHUB_MARKDOWN;
  }

  @Override
  public void generate(Cli cli, PrintWriter writer) {
    this.writer = writer;
    this.listStack = new Stack<String>();

    writer.println("This documentation is generated from " + getClass().getName());
    cli.visit(this);
  }

  @Override
  public void visitFlag(Flag flag) {
  }

  @Override
  public void visitList(List list) {
    writer.println();
    writer.print("# # `");
    writer.print(makeListPrefix());
    writer.print(list.getName());
    writeFlagsText(list.getFlags());
    writer.print("` ...");
    if (list.hasDescription()) {
      writer.print(" : ");
      writer.print(list.getDescription());
    }
    writer.println();
    writeFlagsDescriptions(list.getFlags());
    listStack.push(list.getName());
  }

  @Override
  public void leaveList(List list) {
    listStack.pop();
  }

  @Override
  public void visitCommand(Command command) {
    writer.print("# # #");
    writer.print('`');
    writer.print(makeListPrefix());
    writer.print(command.getName());
    writeFlagsText(command.getFlags());
    writeArgumentsText(command.getArguments());
    writer.print('`');
    writer.println();

    if (command.hasDescription()) {
      writer.print("* ");
      writer.print(command.getDescription());
      writer.println();
    }

    writeFlagsDescriptions(command.getFlags());

    for (Argument argument : command.getArguments()) {
      if (argument.hasDescription()) {
        writer.print("* `");
        writeArgumentText(argument);
        writer.print("` : ");
        writer.print(argument.getDescription());
        writer.println();
      }
    }
  }

  private String makeListPrefix() {
    return (listStack.stream().collect(Collectors.joining(" "))
            + (listStack.empty() ? "" : " "));
  }

  private void writeFlagsDescriptions(Flag[] flags) {
    for (Flag flag : flags) {
      if (flag.hasDescription()) {
        writer.print("* `");
        writeFlagText(flag);
        writer.print("` : ");
        writer.print(flag.getDescription());
        writer.println();
      }
    }
  }

  private void writeFlagsText(Flag[] flags) {
    for (Flag flag : flags) {
      writer.print(' ');
      writeFlagText(flag);
    }
  }

  private void writeFlagText(Flag flag) {
    writer.print('[');
    writer.print(Flag.FLAG_START);
    writer.print(flag.getShortFlag());
    if (flag.hasLongFlag()) {
      writer.print('|');
      writer.print(Flag.FLAG_START);
      writer.print(Flag.FLAG_START);
      writer.print(flag.getLongFlag());
    }
    if (flag.hasArgument()) {
      writer.print(' ');
      writeArgumentText(flag.getArgument());
    }
    writer.print(']');
  }

  private void writeArgumentsText(Argument[] arguments) {
    for (Argument argument : arguments) {
      writer.print(' ');
      writeArgumentText(argument);
    }
  }

  private void writeArgumentText(Argument argument) {
    writer.print('<');
    writer.print(argument.getName());
    writer.print(':');
    writer.print(argument.getType().getName());
    writer.print('>');
  }
}
