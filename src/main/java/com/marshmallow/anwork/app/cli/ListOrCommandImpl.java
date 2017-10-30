package com.marshmallow.anwork.app.cli;

import java.util.ArrayList;

/**
 * This is an element of the CLI tree. Each instance of this class represents a
 * {@link MutableListOrCommand}. Note that this class provides {@link Object#equals(Object)},
 * {@link Object#hashCode()}, and {@link Comparable#compareTo(Object)} implementations for use by
 * derived classes.
 *
 * <p>
 * Created Sep 10, 2017
 * </p>
 *
 * @author Andrew
 */
abstract class ListOrCommandImpl implements MutableListOrCommand, Comparable<ListOrCommandImpl> {

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

    removeFlagForShortFlag(shortFlag);

    MutableFlag flag = new FlagImpl(shortFlag);
    flags.add(flag);
    return flag;
  }

  private void removeFlagForShortFlag(String shortFlag) {
    for (Flag flag : flags) {
      if (flag.getShortFlag().equals(shortFlag)) {
        flags.remove(flag);
        break;
      }
    }
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

  private boolean hasShortFlag(String shortFlag) {
    return findFlagForShortFlag(shortFlag) != null;
  }

  private boolean hasLongFlag(String longFlag) {
    return findShortFlagForLongFlag(longFlag) != null;
  }

  private Flag findFlagForShortFlag(String shortFlag) {
    for (Flag flag : flags) {
      if (flag.getShortFlag().equals(shortFlag)) {
        return flag;
      }
    }
    return null;
  }

  private String findShortFlagForLongFlag(String longFlag) {
    for (Flag flag : flags) {
      if (flag.hasLongFlag() && flag.getLongFlag().equals(longFlag)) {
        return flag.getShortFlag();
      }
    }
    return null;
  }

  /*
   * Section - Parse
   */

  protected boolean isFlag(String arg) {
    return (arg.charAt(0) == Flag.FLAG_START);
  }

  protected int parseFlag(String[] args, int index, ArgumentValues flagValues) {
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
    boolean validFlag = (isLongFlag ? hasLongFlag(flagString) : hasShortFlag(flagString));
    if (!validFlag) {
      throwBadArgException("Unknown flag '" + flagString + "'", args, index);
    }
    String shortFlag = (isLongFlag ? findShortFlagForLongFlag(flagString) : flagString);
    Flag flag = findFlagForShortFlag(shortFlag);

    // Does it have an argument?
    if (flag.hasArgument()) {
      if (index == args.length - 1) {
        throwBadArgException("Expected argument for flag '" + flag + "'", args, index);
      }
      index += 1;
      String argument = args[index];
      Object argumentValue = flag.getArgument().getType().convert(argument);
      flagValues.addValue(shortFlag, argumentValue);
    } else {
      // By default, flags with no argument are set to Boolean.TRUE. See Action#run.
      flagValues.addValue(shortFlag, Boolean.TRUE);
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
