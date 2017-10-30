package com.marshmallow.anwork.app.cli;

import java.util.stream.Stream;

/**
 * This is a {@link ListOrCommandImpl} that represents a {@link List}.
 *
 * <p>
 * Created Oct 15, 2017
 * </p>
 *
 * @author Andrew
 */
class ListImpl extends ListOrCommandImpl implements MutableList {

  private final java.util.List<ListImpl> lists = new java.util.ArrayList<ListImpl>();
  private final java.util.List<CommandImpl> commands = new java.util.ArrayList<CommandImpl>();

  ListImpl(ListImpl parent, String name) {
    super(parent, name);
  }

  @Override
  public MutableList setName(String name) {
    super.setName(name);
    return this;
  }

  @Override
  public MutableList setDescription(String description) {
    super.setDescription(description);
    return this;
  }

  @Override
  public MutableList addList(String name) {
    ListImpl list = new ListImpl(this, name);
    lists.remove(list);
    lists.add(list);
    return list;
  }

  @Override
  public List[] getLists() {
    return lists.toArray(new List[0]);
  }

  @Override
  public MutableCommand addCommand(String name, Action action) {
    CommandImpl command = new CommandImpl(this, name, action);
    commands.remove(command);
    commands.add(command);
    return command;
  }

  @Override
  public Command[] getCommands() {
    return commands.toArray(new Command[0]);
  }

  /*
   * Section - Parsing
   */

  void parse(String[] args) {
    ArgumentValues flagValues = new ArgumentValues();
    if (!parseList(args, 0, flagValues)) {
      printUsage();
    }
  }

  // Returns true iff a command was run.
  private boolean parseList(String[] args, int index, ArgumentValues flagValues) {
    while (index < args.length) {
      String arg = args[index];
      if (isFlag(arg)) {
        index = parseFlag(args, index, flagValues);
      } else if (hasList(arg)) {
        return findList(arg).parseList(args, index + 1, flagValues);
      } else if (hasCommand(arg)) {
        findCommand(arg).parse(args, index + 1, flagValues);
        return true;
      } else {
        // This is not a list or a command, so let's call it an unknown command!
        throwBadArgException("Unknown command/list '" + arg + "' for list " + getName(),
                             args,
                             index);
      }
    }
    return false;
  }

  private boolean hasList(String name) {
    return findList(name) != null;
  }

  private boolean hasCommand(String name) {
    return findCommand(name) != null;
  }

  private ListImpl findList(String name) {
    for (ListImpl list : lists) {
      if (list.getName().equals(name)) {
        return list;
      }
    }
    return null;
  }

  private CommandImpl findCommand(String name) {
    for (CommandImpl command : commands) {
      if (command.getName().equals(name)) {
        return command;
      }
    }
    return null;
  }

  private void printUsage() {
    StringBuilder builder = new StringBuilder();
    printUsage(this, builder);
    System.out.println(builder.toString());
  }

  private void printUsage(List list, StringBuilder builder) {
    TextDocumentationGenerator.writeList(new List[] { list }, builder);
    if (list.getFlags().length > 0) {
      TextDocumentationGenerator.writeFlags(list.getFlags(), builder);
    }
    if (list.getCommands().length > 0) {
      TextDocumentationGenerator.writeCommands(list.getCommands(), builder);
    }
    Stream.of(list.getLists()).forEach((l) -> printUsage(l, builder));
  }

  /*
   * Section - Visitor
   */

  void visit(Visitor visitor) {
    visitor.visitList(this);
    visitFlags(visitor);
    commands.stream().sorted().forEach((command) -> command.visit(visitor));
    lists.stream().sorted().forEach((list) -> list.visit(visitor));
    visitor.leaveList(this);
  }
}