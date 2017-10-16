package com.marshmallow.anwork.app.cli;

import java.io.PrintWriter;
import java.io.StringWriter;
import java.util.ArrayList;
import java.util.Stack;
import java.util.stream.Collectors;
import java.util.stream.Stream;

/**
 * This is a {@link DocumentationGenerator} that is useful for simple text formats, like on the
 * command line.
 *
 * <p>The format for this documentation looks something like this.
 * <pre>
 * root-list ... : Description for root-list
 *   Flags
 *     -d|--debug                  : Description for debug flag
 *     -f|--file (filename:STRING) : Description for file flag
 *   Commands
 *     command-1      : Description for command-1
 *     command-2      : Description for command-2
 *     long-command-3 : Description for long-command-3
 *
 * root-list list-a ... : Description for root-list
 *   Commands
 *     command-abc : Description for command-abc
 *     foo         : Description for foo
 *
 * root-list list-b ... : Description for list-b
 *   Commands
 *     bar      : Description for bar command
 *     bat-tuna : Description for bat-tuna command
 * </pre>
 *
 * <p>
 * Created Oct 16, 2017
 * </p>
 *
 * @author Andrew
 */
class TextDocumentationGenerator implements DocumentationGenerator, Visitor {

  /**
   * This type holds per-{@link List} state when generating documentation.
   *
   * <p>
   * Created Oct 16, 2017
   * </p>
   *
   * @author Andrew
   */
  private static class ListState {
    private final StringBuilder textBuilder = new StringBuilder();
    private final List list;
    private final java.util.List<Flag> flags = new ArrayList<Flag>();
    private final java.util.List<Command> commands = new ArrayList<Command>();

    public ListState(List list) {
      this.list = list;
    }

    public StringBuilder getTextBuilder() {
      return textBuilder;
    }

    public List getList() {
      return list;
    }

    public void addFlag(Flag flag) {
      flags.add(flag);
    }

    public Flag[] getFlags() {
      return flags.toArray(new Flag[0]);
    }

    public void addCommand(Command command) {
      commands.add(command);
    }

    public Command[] getCommands() {
      return commands.toArray(new Command[0]);
    }
  }

  private static final String FLAGS = "  Flags\n";
  private static final String COMMANDS = "  Commands\n";

  private Stack<ListState> listStateStack = new Stack<ListState>();

  @Override
  public DocumentationType getType() {
    return DocumentationType.TEXT;
  }

  @Override
  public void generate(Cli cli, PrintWriter writer) {
    cli.visit(this);
    while (!listStateStack.isEmpty()) {
      writer.print(listStateStack.pop().getTextBuilder().toString());
    }
  }

  @Override
  public void visitFlag(Flag flag) {
    listStateStack.peek().addFlag(flag);
  }

  @Override
  public void visitList(List list) {
    listStateStack.push(new ListState(list));
  }

  @Override
  public void leaveList(List list) {
    writeList();
    writeFlags();
    writeCommands();
    listStateStack.pop();
  }

  @Override
  public void visitCommand(Command command) {
    listStateStack.peek().addCommand(command);
  }

  private void writeList() {
    ListState state = listStateStack.peek();
    if (listStateStack.size() > 1) {
      state.getTextBuilder().append('\n');
    }
    List[] lists = listStateStack.stream().map(ListState::getList).toArray(List[]::new);
    writeList(lists, state.getTextBuilder());
  }

  /**
   * Write out an array of {@link List}'s in (hopefully) a readable way.
   *
   * @param lists The {@link List}'s to write out
   * @param builder The {@link StringBuilder} to use to write the {@link List}'s
   */
  static void writeList(List[] lists, StringBuilder builder) {
    List currentList = lists[lists.length - 1];
    String listsName = Stream.of(lists).map(List::getName).collect(Collectors.joining(" "));
    builder.append(listsName);
    builder.append(" ...");
    if (currentList.hasDescription()) {
      writeDescription(currentList.getDescription(), builder);
    }
    builder.append('\n');
  }

  private void writeFlags() {
    ListState state = listStateStack.peek();
    Flag[] flags = state.getFlags();
    if (flags.length > 0) {
      writeFlags(flags, state.getTextBuilder());
    }
  }

  /**
   * Write out an array of {@link Flag}'s in (hopefully) a readable way.
   *
   * @param flags The {@link Flag}'s to write out
   * @param builder The {@link StringBuilder} to use to write the {@link Flag}'s
   */
  static void writeFlags(Flag[] flags, StringBuilder builder) {
    int descriptionColonIndent = getLongestFlagName(flags);
    builder.append(FLAGS);
    for (Flag flag : flags) {
      String format = "%-" + descriptionColonIndent + "s";
      builder.append(String.format(format, getFlagDisplayName(flag)));
      if (flag.hasDescription()) {
        writeDescription(flag.getDescription(), builder);
      }
      builder.append('\n');
    }
  }

  private static int getLongestFlagName(Flag[] flags) {
    int longest = 0;
    for (Flag flag : flags) {
      int flagNameLength = getFlagDisplayName(flag).length();
      if (flagNameLength > longest) {
        longest = flagNameLength;
      }
    }
    return longest;
  }

  private static String getFlagDisplayName(Flag flag) {
    StringBuilder builder = new StringBuilder();
    builder.append("    "); // indent
    builder.append(Flag.FLAG_START).append(flag.getShortFlag());
    if (flag.hasLongFlag()) {
      builder.append('|')
             .append(Flag.FLAG_START)
             .append(Flag.FLAG_START)
             .append(flag.getLongFlag());
    }
    if (flag.hasArgument()) {
      Argument argument = flag.getArgument();
      builder.append(' ')
             .append('(')
             .append(argument.getName())
             .append(':')
             .append(argument.getType().getName())
             .append(')');
    }
    return builder.toString();
  }

  private void writeCommands() {
    ListState state = listStateStack.peek();
    Command[] commands = state.getCommands();
    if (commands.length > 0) {
      writeCommands(commands, state.getTextBuilder());
    }
  }

  /**
   * Write out an array of {@link Command}'s in (hopefully) a readable way.
   *
   * @param commands The {@link Command}'s to write out
   * @param builder The {@link StringBuilder} to use to write the {@link Command}'s
   */
  static void writeCommands(Command[] commands, StringBuilder builder) {
    int longestCommandName = getLongestCommandName(commands);
    builder.append(COMMANDS);
    for (Command command : commands) {
      String format = "    %-" + longestCommandName + "s";
      builder.append(String.format(format, command.getName()));
      if (command.hasDescription()) {
        writeDescription(command.getDescription(), builder);
      }
      builder.append('\n');
    }
  }

  private static int getLongestCommandName(Command[] commands) {
    return Stream.of(commands)
                 .map(command -> command.getName().length())
                 .max((a, b) -> a - b)
                 .get();
  }

  private static void writeDescription(String description, StringBuilder builder) {
    builder.append(" : ");
    builder.append(description);
  }
}
