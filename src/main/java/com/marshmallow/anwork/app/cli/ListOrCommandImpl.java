package com.marshmallow.anwork.app.cli;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * This is an element of the CLI tree. Each instance of this class represents a
 * {@link MutableListOrCommand}.
 *
 * <p>
 * Created Sep 10, 2017
 * </p>
 *
 * @author Andrew
 */
abstract class ListOrCommandImpl implements MutableListOrCommand, Comparable<ListOrCommandImpl> {

  /**
   * This is a helper class for CLI parsing functionality. It contains information about the per
   * {@link #parse(String[])} state. Calling {@link #setActiveNode(ListOrCommandImpl)} will initialize the
   * state for the current {@link ListOrCommandImpl} being {@link ListOrCommandImpl#parse(String[])}'d.
   *
   * <p>
   * Created Sep 10, 2017
   * </p>
   *
   * @author Andrew
   */
  private static class ParseContext {

    private ListOrCommandImpl activeNode;
    private final List<String> parameters = new ArrayList<String>();
    private final ArgumentValues flagValues = new ArgumentValues();

    // This is a map from Flag#getShortFlag to Flag
    private final Map<String, Flag> flags = new LinkedHashMap<String, Flag>();
    // This is a map from Flag#getLongFlag to Flag#getShortFlag
    private final Map<String, String> longFlags = new LinkedHashMap<String, String>();
    // This is a map from Node#getName to Node
    private final Map<String, ListOrCommandImpl> children = new LinkedHashMap<String, ListOrCommandImpl>();

    private void reinitialize() {
      flags.clear();
      for (Flag flag : activeNode.flags) {
        flags.put(flag.getShortFlag(), flag);
        if (flag.hasLongFlag()) {
          longFlags.put(flag.getLongFlag(), flag.getShortFlag());
        }
      }

      children.clear();
      for (ListOrCommandImpl child : activeNode.children) {
        children.put(child.getName(), child);
      }
    }

    public ListOrCommandImpl getActiveNode() {
      return activeNode;
    }

    public void setActiveNode(ListOrCommandImpl activeNode) {
      this.activeNode = activeNode;
      reinitialize();
    }

    public String[] getParameters() {
      return parameters.toArray(new String[0]);
    }

    public void addParameter(String parameter) {
      parameters.add(parameter);
    }

    public ArgumentValues getFlagValues() {
      return flagValues;
    }

    public boolean hasFlag(String shortFlag) {
      return flags.containsKey(shortFlag);
    }

    public Flag getFlag(String shortFlag) {
      return flags.get(shortFlag);
    }

    public boolean hasLongFlag(String longFlag) {
      return longFlags.containsKey(longFlag);
    }

    public String getShortFlag(String shortFlag) {
      return longFlags.get(shortFlag);
    }

    public boolean hasChild(String name) {
      return children.containsKey(name);
    }

    public ListOrCommandImpl getChild(String name) {
      return children.get(name);
    }
  }

  private final java.util.List<Flag> flags = new ArrayList<Flag>();
  private final java.util.List<ListOrCommandImpl> children = new ArrayList<ListOrCommandImpl>();

  private String name;
  private Action action;
  private String description;

  protected ListOrCommandImpl(String name, Action action) {
    this.name = name;
    this.action = action;
  }

  /**
   * Return whether or not this {@link ListOrCommandImpl} represents a CLI {@link List}.
   *
   * @return Whether or not this {@link ListOrCommandImpl} represents a CLI {@link List}
   */
  protected abstract boolean isList();

  @Override
  public String toString() {
    return String.format("%s:%s", getClass().getName(), getName());
  }

  @Override
  public boolean equals(Object other) {
    if (other == null) {
      return this == null;
    } else if (!(other instanceof ListOrCommandImpl)) {
      return false;
    } else {
      return ((ListOrCommandImpl)other).getName().equals(getName());
    }
  }

  @Override
  public int hashCode() {
    return getName().hashCode();
  }

  @Override
  public MutableListOrCommand setName(String name) {
    this.name = name;
    return this;
  }

  @Override
  public String getName() {
    return name;
  }

  @Override
  public MutableListOrCommand setDescription(String description) {
    this.description = description;
    return this;
  }

  @Override
  public boolean hasDescription() {
    return description != null;
  }

  @Override
  public String getDescription() {
    return description;
  }

  protected MutableListOrCommand setAction(Action action) {
    this.action = action;
    return this;
  }

  protected Action getAction() {
    return action;
  }

  /*
   * Section - Flags
   */

  @Override
  public Flag[] getFlags() {
    return flags.toArray(new Flag[0]);
  }

  @Override
  public MutableFlag addFlag(String shortFlag) {
    for (Flag flag : flags) {
      if (flag.getShortFlag().equals(shortFlag)) {
        flags.remove(flag);
        break;
      }
    }

    MutableFlag flag = new FlagImpl(shortFlag);
    flags.add(flag);
    return flag;
  }

  /*
   * Section - Children
   */

  protected void addChild(ListOrCommandImpl child) {
    if (children.contains(child)) {
      children.remove(child);
    }
    children.add(child);
  }

  protected ListOrCommandImpl[] getChildren(boolean lists) {
    return children.stream()
                   .filter(lists ? (node) -> node.isList() : (node) -> !node.isList())
                   .toArray(ListOrCommandImpl[]::new);
  }

  /*
   * Section - Usage
   */

  String getUsage() {
    return makeCommandUsage(makeFlagUsage());
  }

  private String makeFlagUsage() {
    StringBuilder builder = new StringBuilder();
    for (Flag flag : flags) {
      builder.append('[');
      builder.append(Flag.FLAG_START).append(flag.getShortFlag());
      if (flag.hasLongFlag()) {
        builder.append('|').append(Flag.FLAG_START).append(Flag.FLAG_START);
        builder.append(flag.getLongFlag());
      }
      if (flag.hasArgument()) {
        builder.append(' ').append('<').append(flag.getArgument().getName()).append('>');
      }
      builder.append(' ').append(flag.getDescription());
      builder.append(']');
      builder.append(' ');
    }
    return builder.toString();
  }

  private String makeCommandUsage(String flagUsage) {
    if (children.size() == 0) {
      return flagUsage;
    }

    StringBuilder builder = new StringBuilder();
    for (ListOrCommandImpl child : children) {
      builder.append(name).append(' ');
      builder.append(flagUsage);
      builder.append(child.name).append(" : ").append(child.description);
      builder.append('\n');
    }
    return builder.toString();
  }

  /*
   * Section - Parse
   */

  void parse(String[] args) {
    ParseContext context = new ParseContext();
    parse(args, 0, context);
    validateContext(context);
    runActiveNodeFromContext(context);
  }

  private int parse(String[] args, int index, ParseContext context) {
    context.setActiveNode(this);
    while (index < args.length) {
      index = parseArg(args, index, context);
    }
    return index;
  }

  // Process the argument at index. Returns the next index to process.
  private int parseArg(String[] args, int index, ParseContext context) {
    String arg = args[index];
    int nextIndex;
    if (isFlag(arg)) {
      nextIndex = parseFlag(args, index, context);
    } else if (context.hasChild(arg) && context.getParameters().length == 0) {
      ListOrCommandImpl child = context.getChild(arg);
      nextIndex = child.parse(args, index + 1, context);
    } else {
      context.addParameter(arg);
      nextIndex = index + 1;
    }
    return nextIndex;
  }

  private boolean isFlag(String arg) {
    return (arg.charAt(0) == Flag.FLAG_START);
  }

  private int parseFlag(String[] args, int index, ParseContext context) {
    String arg = args[index];

    // Is it valid flag syntax?
    boolean isLongFlag = false;
    if (arg.charAt(0) != Flag.FLAG_START) {
      throwBadArgException("Expected flag syntax", args, index);
    }
    if (arg.charAt(1) == Flag.FLAG_START) {
      isLongFlag = true;
    }

    // Is it a valid flag?
    String flagString = arg.substring(isLongFlag ? 2 : 1);
    boolean validFlag = (isLongFlag
                         ? context.hasLongFlag(flagString)
                         : context.hasFlag(flagString));
    if (!validFlag) {
      throwBadArgException("Unknown flag '" + flagString + "'", args, index);
    }
    String shortFlag = (isLongFlag ? context.getShortFlag(flagString) : flagString);
    Flag flag = context.getFlag(shortFlag);

    // Does it have an argument?
    if (flag.hasArgument()) {
      if (index == args.length - 1) {
        throwBadArgException("Expected argument for flag '" + flag + "'", args, index);
      }
      index += 1;
      String argument = args[index];
      Object argumentValue = flag.getArgument().getType().convert(argument);
      context.getFlagValues().addShortFlagValue(shortFlag, argumentValue);
    } else {
      // By default, flags with no argument are set to Boolean.TRUE. See Action#run.
      context.getFlagValues().addShortFlagValue(shortFlag, Boolean.TRUE);
    }

    return index + 1;
  }

  private void throwBadArgException(String baseMessage, String[] args, int index) {
    String message = String.format("%s index=%d, arg=%s", baseMessage, index, args[index]);
    throw new IllegalArgumentException(message);
  }

  private void validateContext(ParseContext context) {
    ListOrCommandImpl activeNode = context.getActiveNode();
    String[] parameters = context.getParameters();
    if (activeNode.isList() && parameters.length != 0) {
      throw new IllegalArgumentException("Unknown command '" + parameters[0]
                                         + "' for list " + activeNode.name);
    }
  }

  private void runActiveNodeFromContext(ParseContext context) {
    ListOrCommandImpl activeNode = context.getActiveNode();
    String[] parameters = context.getParameters();
    activeNode.action.run(context.getFlagValues(), parameters);
  }

  /*
   * Section - Visitor
   */

  /**
   * Start the visitation of a {@link ListOrCommandImpl}.
   *
   * @param visitor The visitor that is visiting during this visitation
   */
  protected abstract void startVisit(Visitor visitor);

  /**
   * End the visitation of a {@link ListOrCommandImpl}.
   *
   * @param visitor The visitor that is visiting during this visitation
   */
  protected abstract void endVisit(Visitor visitor);

  void visit(Visitor visitor) {
    // First, we visit ourselves.
    startVisit(visitor);
    // Second, we visit our flags.
    flags.stream()
         .sorted()
         .forEach(flag -> visitor.visitFlag(flag));
    // Third, we visit our commands.
    children.stream()
            .filter(node -> !node.isList())
            .sorted()
            .forEach(command -> command.visit(visitor));
    // Fourth, we visit our lists.
    children.stream()
            .filter(node -> node.isList())
            .sorted()
            .forEach(list -> list.visit(visitor));
    endVisit(visitor);
  }

  @Override
  public int compareTo(ListOrCommandImpl other) {
    return name.compareTo(other.name);
  }
}
