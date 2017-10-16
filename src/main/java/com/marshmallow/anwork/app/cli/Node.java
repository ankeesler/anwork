package com.marshmallow.anwork.app.cli;

import java.util.HashMap;
import java.util.LinkedHashMap;
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
abstract class Node implements MutableListOrCommand, Comparable<Node> {

  // This is a map from short flag to information about the flag.
  private Map<String, Flag> shortFlagInfo = new LinkedHashMap<String, Flag>();
  // This is a map from long flag to short flag.
  private Map<String, String> longFlagShortFlags = new HashMap<String, String>();
  // This is a map from child name to child node.
  private Map<String, Node> children = new HashMap<String, Node>();

  private String name;
  private Action action;
  private String description = "";

  protected Node(String name, Action action) {
    this.name = name;
    this.action = action;
  }

  protected abstract boolean isList();

  @Override
  public String toString() {
    return String.format("%s:%s", getClass().getName(), getName());
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
    return shortFlagInfo.values().toArray(new Flag[0]);
  }

  @Override
  public MutableFlag addFlag(String shortFlag) {
    MutableFlag flag = new FlagImpl(shortFlag);
    shortFlagInfo.put(flag.getShortFlag(), flag);
    return flag;
  }

  /*
   * Section - Children
   */

  protected void addChild(Node child) {
    children.put(child.name, child);
  }

  protected Node[] getChildren(boolean lists) {
    return children.values()
                   .stream()
                   .filter(lists ? (node) -> node.isList() : (node) -> !node.isList())
                   .toArray(Node[]::new);
  }

  /*
   * Section - Usage
   */

  String getUsage() {
    return makeCommandUsage(makeFlagUsage());
  }

  private String makeFlagUsage() {
    StringBuilder builder = new StringBuilder();
    for (Flag flag : shortFlagInfo.values()) {
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
    for (Node child : children.values()) {
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
    updateLongFlagsShortFlags();
    context.setActiveNode(this);
    while (index < args.length) {
      index = parseArg(args, index, context);
    }
    return index;
  }

  private void updateLongFlagsShortFlags() {
    // We do this lazily so that the flags can be updated until the last moment when #parse is
    // called.
    longFlagShortFlags.clear();
    for (Flag flag : shortFlagInfo.values()) {
      if (flag.hasLongFlag()) {
        longFlagShortFlags.put(flag.getLongFlag(), flag.getShortFlag());
      }
    }
  }

  // Process the argument at index. Returns the next index to process.
  private int parseArg(String[] args, int index, ParseContext context) {
    String arg = args[index];
    int nextIndex;
    if (isFlag(arg)) {
      nextIndex = parseFlag(args, index, context);
    } else if (isChild(arg) && context.getParameters().length == 0) {
      Node child = getChild(arg);
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
                         ? longFlagShortFlags.containsKey(flagString)
                         : shortFlagInfo.containsKey(flagString));
    if (!validFlag) {
      throwBadArgException("Unknown flag '" + flagString + "'", args, index);
    }
    String shortFlag = (isLongFlag ? longFlagShortFlags.get(flagString) : flagString);
    Flag flag = shortFlagInfo.get(shortFlag);

    // Does it have an argument?
    if (flag.hasArgument()) {
      if (index == args.length - 1) {
        throwBadArgException("Expected argument for flag '" + flag + "'", args, index);
      }
      index += 1;
      String argument = args[index];
      Object argumentValue = flag.getArgument().getType().convert(argument);
      context.getFlags().addShortFlagValue(shortFlag, argumentValue);
    } else {
      // By default, flags with no argument are set to Boolean.TRUE.
      context.getFlags().addShortFlagValue(shortFlag, Boolean.TRUE);
    }

    return index + 1;
  }

  private boolean isChild(String arg) {
    return children.containsKey(arg);
  }

  private Node getChild(String arg) {
    return children.get(arg);
  }

  private void throwBadArgException(String baseMessage, String[] args, int index) {
    String message = String.format("%s index=%d, arg=%s", baseMessage, index, args[index]);
    throw new IllegalArgumentException(message);
  }

  private void validateContext(ParseContext context) {
    Node activeNode = context.getActiveNode();
    String[] parameters = context.getParameters();
    if (activeNode.isList() && parameters.length != 0) {
      throw new IllegalArgumentException("Unknown command '" + parameters[0]
                                         + "' for list " + activeNode.name);
    }
  }

  private void runActiveNodeFromContext(ParseContext context) {
    Node activeNode = context.getActiveNode();
    String[] parameters = context.getParameters();
    activeNode.action.run(context.getFlags(), parameters);
  }

  /*
   * Section - Visitor
   */

  protected abstract void startVisit(Visitor visitor);

  protected abstract void endVisit(Visitor visitor);

  void visit(Visitor visitor) {
    // First, we visit ourselves.
    startVisit(visitor);
    // Second, we visit our flags.
    shortFlagInfo.values()
                 .stream()
                 .sorted()
                 .forEach(flag -> visitor.visitFlag(flag));
    // Third, we visit our commands.
    children.values()
            .stream()
            .filter(node -> !node.isList())
            .sorted()
            .forEach(command -> command.visit(visitor));
    // Fourth, we visit our lists.
    children.values()
            .stream()
            .filter(node -> node.isList())
            .sorted()
            .forEach(list -> list.visit(visitor));
    endVisit(visitor);
  }

  @Override
  public int compareTo(Node other) {
    return name.compareTo(other.name);
  }
}
