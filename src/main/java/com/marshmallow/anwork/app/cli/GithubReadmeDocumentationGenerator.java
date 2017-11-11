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

    StringBuilder lineBuilder = new StringBuilder();
    lineBuilder.append('`');
    lineBuilder.append(makeListPrefix());
    lineBuilder.append(list.getName());
    lineBuilder.append(makeFlagText(list.getFlags()));
    lineBuilder.append('`');
    lineBuilder.append(" ...");
    if (list.hasDescription()) {
      lineBuilder.append(" : ");
      lineBuilder.append(list.getDescription());
    }
    writer.println(lineBuilder.toString());
    listStack.push(list.getName());
  }

  @Override
  public void leaveList(List list) {
  }

  @Override
  public void visitCommand(Command command) {
    StringBuilder lineBuilder = new StringBuilder();
    lineBuilder.append("* ");
    lineBuilder.append('`');
    lineBuilder.append(makeListPrefix());
    lineBuilder.append(command.getName());
    lineBuilder.append(makeFlagText(command.getFlags()));
    for (Argument argument : command.getArguments()) {
      lineBuilder.append(makeArgumentText(argument));
    }
    lineBuilder.append('`');
    writer.println(lineBuilder.toString());

    if (command.hasDescription()) {
      writer.println("* * " + command.getDescription());
    }
  }

  private String makeListPrefix() {
    return (listStack.stream().collect(Collectors.joining(" "))
            + (listStack.empty() ? "" : " "));
  }

  private String makeFlagText(Flag[] flags) {
    StringBuilder textBuilder = new StringBuilder();
    for (Flag flag : flags) {
      textBuilder.append(' ');
      textBuilder.append('[');
      textBuilder.append(Flag.FLAG_START).append(flag.getShortFlag());
      if (flag.hasLongFlag()) {
        textBuilder.append('|')
                   .append(Flag.FLAG_START)
                   .append(Flag.FLAG_START)
                   .append(flag.getLongFlag());
      }
      if (flag.hasArgument()) {
        textBuilder.append(makeArgumentText(flag.getArgument()));
      }
      textBuilder.append(']');
    }
    return textBuilder.toString();
  }

  private static String makeArgumentText(Argument argument) {
    return new StringBuilder().append(' ')
                              .append('<')
                              .append(argument.getName())
                              .append(':')
                              .append(argument.getType().getName())
                              .append('>')
                              .toString();
  }
}
