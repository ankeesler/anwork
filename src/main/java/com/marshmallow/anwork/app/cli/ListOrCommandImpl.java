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
   * {@link #parse(String[])} state. Calling {@link #setActiveNode(ListOrCommandImpl)} will
   * initialize the state for the current {@link ListOrCommandImpl} being
   * {@link #parse(String[])}'d.
   *
   * <p>
   * Created Sep 10, 2017
   * </p>
   *
   * @author Andrew
   */
  static class ParseContext {

    private final List<String> arguments = new ArrayList<String>();
    private final ArgumentValues flagValues = new ArgumentValues();

    // This is a map from Flag#getShortFlag to Flag
    private final Map<String, Flag> flags
        = new LinkedHashMap<String, Flag>();
    // This is a map from Flag#getLongFlag to Flag#getShortFlag
    private final Map<String, String> longFlags
        = new LinkedHashMap<String, String>();

    private void reinitialize(ListOrCommand activeNode) {
      flags.clear();
      for (Flag flag : activeNode.getFlags()) {
        flags.put(flag.getShortFlag(), flag);
        if (flag.hasLongFlag()) {
          longFlags.put(flag.getLongFlag(), flag.getShortFlag());
        }
      }
    }

    public void setActiveNode(ListOrCommand activeNode) {
      reinitialize(activeNode);
    }

    public String[] getArguments() {
      return arguments.toArray(new String[0]);
    }

    public void addArgument(String argument) {
      arguments.add(argument);
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
  }

  private final ListOrCommandImpl parent;
  private String name;
  private String description;
  private final java.util.List<Flag> flags;

  protected ListOrCommandImpl(ListOrCommandImpl parent, String name) {
    this.parent = parent;
    this.name = name;
    flags = new ArrayList<Flag>();
  }

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
  public int compareTo(ListOrCommandImpl other) {
    return name.compareTo(other.name);
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

  /*
   * Section - Flags
   */

  @Override
  public Flag[] getFlags() {
    return flags.toArray(new Flag[0]);
  }

  @Override
  public MutableFlag addFlag(String shortFlag) {
    if (parentContainsFlag(shortFlag)) {
      throw new IllegalStateException("Cannot assign the same short flag "
                                      + "to hiearchical lists/commands! "
                                      + "If this was allowed, then there would be ambiguity upon "
                                      + "accessing whether or not a flag was set on a command "
                                      + "when the same flag was set on a parent list.");
    }

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

  private boolean parentContainsFlag(String shortFlag) {
    if (parent == null) {
      return false;
    }

    for (Flag flag : parent.getFlags()) {
      if (flag.getShortFlag().equals(shortFlag)) {
        return true;
      }
    }

    return parent.parentContainsFlag(shortFlag);
  }

  /*
   * Section - Parse
   */

  protected boolean isFlag(String arg) {
    return (arg.charAt(0) == Flag.FLAG_START);
  }

  protected int parseFlag(String[] args, int index, ParseContext context) {
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
      context.getFlagValues().addValue(shortFlag, argumentValue);
    } else {
      // By default, flags with no argument are set to Boolean.TRUE. See Action#run.
      context.getFlagValues().addValue(shortFlag, Boolean.TRUE);
    }

    return index + 1;
  }

  protected void throwBadArgException(String baseMessage, String[] args, int index) {
    String message = String.format("%s index=%d, arg=%s", baseMessage, index, args[index]);
    throw new IllegalArgumentException(message);
  }

  /*
   * Section - Visitor
   */

  protected void visitFlags(Visitor visitor) {
    flags.stream().sorted().forEach((flag) -> visitor.visitFlag(flag));
  }
}
