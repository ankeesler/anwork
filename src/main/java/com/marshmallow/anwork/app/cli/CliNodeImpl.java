package com.marshmallow.anwork.app.cli;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * This is an implementation of a CLI node.
 *
 * This class is not meant to be instantiated outside this package. Users should
 * interface through {@link Cli} and call {@link CliNode#addCommand(String)).
 *
 * @author Andrew
 * @date Sep 9, 2017
 */
class CliNodeImpl implements CliNode {

  // This is a map from short flag to information about the flag.
  private Map<String, CliFlag> shortFlagInfo = new LinkedHashMap<String, CliFlag>();
  // This is a map from long flag to short flag.
  private Map<String, String> longFlagShortFlags = new HashMap<String, String>();
  // This is a map from child name to child node.
  private Map<String, CliNode> children = new HashMap<String, CliNode>();

  private String name;
  private String description;
  private CliAction action;

  CliNodeImpl(String name, String description, CliAction action) {
    this.name = name;
    this.description = description;
    this.action = action;
  }

  @Override
  public void addFlag(CliFlag flag) {
    shortFlagInfo.put(flag.getShortFlag(), flag);
    if (flag.hasLongFlag()) {
      longFlagShortFlags.put(flag.getLongFlag(), flag.getShortFlag());
    }
  }

  @Override
  public CliNode addCommand(String name, String description, CliAction action) {
    CliNode child = new CliNodeImpl(name, description, action);
    children.put(name, child);
    return child;
  }

  @Override
  public String getName() {
    return name;
  }

  @Override
  public String getDescription() {
    return description;
  }

  @Override
  public String getUsage() {
    return makeCommandUsage(makeFlagUsage());
  }

  private String makeFlagUsage() {
    StringBuilder builder = new StringBuilder();
    for (CliFlag flag : shortFlagInfo.values()) {
      builder.append(CliFlag.FLAG_START).append(flag.getShortFlag());
      if (flag.hasLongFlag()) {
        builder.append('|').append(CliFlag.FLAG_START).append(CliFlag.FLAG_START).append(flag.getLongFlag());
      }
      if (flag.hasParameter()) {
        builder.append(' ').append('<').append(flag.getParameterName()).append('>');
      }
      builder.append(' ').append(flag.getDescription());
      builder.append('\n');      
    }
    return builder.toString();
  }

  private String makeCommandUsage(String flagUsage) {
    if (children.size() == 0) {
      return flagUsage;
    }

    StringBuilder builder = new StringBuilder();
    for (CliNode child : children.values()) {
      builder.append(name).append(' ');
      builder.append(flagUsage).append(' ');
      builder.append(child.getName()).append(" : ").append(child.getDescription());
      builder.append('\n');
    }
    return builder.toString();
  }

  @Override
  public void parse(String[] args) {
    int i = 0;
    while (i < args.length) {
      i = processArg(args, i);
    }
  }

  // Process the argument at index. Returns the next index to process.
  private int processArg(String[] args, int index) throws IllegalArgumentException {
    String arg = args[index];

    // Is it valid flag syntax?
    boolean isLongFlag = false;
    if (arg.charAt(0) != CliFlag.FLAG_START) {
      throwBadArgException("Expected flag syntax", args, index);
    }
    if (arg.charAt(1) == CliFlag.FLAG_START) {
      isLongFlag = true;
    }

    // Is it a valid flag?
    String flagString = arg.substring(isLongFlag ? 2 : 1);
    boolean validFlag = (isLongFlag ? longFlagShortFlags.containsKey(flagString) : shortFlagInfo.containsKey(flagString));
    if (!validFlag) {
      throwBadArgException("Unknown flag '" + flagString + "'", args, index);
    }
    String shortFlag = (isLongFlag ? longFlagShortFlags.get(flagString) : flagString);
    CliFlag flag = shortFlagInfo.get(shortFlag);

    // Does it have an argument?
    List<String> arguments = new ArrayList<String>();
    if (flag.hasParameter()) {
      if (index == args.length -1) {
        throwBadArgException("Expected argument for flag '" + flag + "'", args, index);
      }
      index += 1;
      arguments.add(args[index]);
    }

    flag.getAction().run(arguments.toArray(new String[0]));

    return index + 1;
  }

  private void throwBadArgException(String baseMessage, String args[], int index) throws IllegalArgumentException {
    String message = String.format("%s index=%d, arg=%s", baseMessage, index, args[index]);
    throw new IllegalArgumentException(message);
  }
}
