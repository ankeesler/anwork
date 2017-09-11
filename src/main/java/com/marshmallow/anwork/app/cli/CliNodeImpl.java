package com.marshmallow.anwork.app.cli;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * This is an element of the CLI tree.
 *
 * @author Andrew
 * Created Sep 10, 2017
 */
class CliNodeImpl implements CliList, CliCommand {

  private static final int LIST_PARAM_COUNT = -1;

  // Make a root node.
  static CliNodeImpl makeRoot(String name, String description) {
    return new CliNodeImpl(name, description, LIST_PARAM_COUNT, null);
  }

  // This is a map from short flag to information about the flag.
  private Map<String, CliFlag> shortFlagInfo = new LinkedHashMap<String, CliFlag>();
  // This is a map from long flag to short flag.
  private Map<String, String> longFlagShortFlags = new HashMap<String, String>();
  // This is a map from child name to child node.
  private Map<String, CliNodeImpl> children = new HashMap<String, CliNodeImpl>();

  private final String name;
  private final String description;
  private final int paramCount;
  private final CliAction action;

  // See static #make... methods above.
  private CliNodeImpl(String name, String description, int paramCount, CliAction action) {
    this.name = name;
    this.description = description;
    this.paramCount = paramCount;
    if (action == null) {
      this.action = new CliUsageAction(this);
    } else {
      this.action = action;
    }
  }

  private boolean isList() {
    return paramCount == LIST_PARAM_COUNT;
  }

  /*
   * Section - Flags
   */

  @Override
  public void addShortFlag(String shortFlag, String description, CliAction action) {
    addFlag(CliFlag.makeShortFlag(shortFlag, description, action));
  }

  @Override
  public void addShortFlagWithParameter(String shortFlag,
                                        String description,
                                        String parameterName,
                                        CliAction action) {
    addFlag(CliFlag.makeShortFlagWithParameter(shortFlag, description, parameterName, action));
  }

  @Override
  public void addLongFlag(String shortFlag,
                          String longFlag,
                          String description,
                          CliAction action) {
    addFlag(CliFlag.makeLongFlag(shortFlag, longFlag, description, action));
  }

  @Override
  public void addLongFlagWithParameter(String shortFlag,
                                       String longFlag,
                                       String description,
                                       String parameterName,
                                       CliAction action) {
    addFlag(CliFlag.makeLongFlagWithParameter(shortFlag,
                                              longFlag,
                                              description,
                                              parameterName,
                                              action));
  }

  private void addFlag(CliFlag flag) {
    shortFlagInfo.put(flag.getShortFlag(), flag);
    if (flag.hasLongFlag()) {
      longFlagShortFlags.put(flag.getLongFlag(), flag.getShortFlag());
    }
  }

  /*
   * Section - Children
   */

  @Override
  public CliNodeImpl addList(String name, String description) {
    CliNodeImpl list = new CliNodeImpl(name, description, LIST_PARAM_COUNT, null);
    addChild(list);
    return list;
  }

  @Override
  public CliNodeImpl addCommand(String name, String description, CliAction action) {
    CliNodeImpl command = new CliNodeImpl(name, description, 0, action);
    addChild(command);
    return command;
  }

  private void addChild(CliNodeImpl child) {
    children.put(child.name, child);
  }

  /*
   * Section - Usage
   */

  String getUsage() {
    return makeCommandUsage(makeFlagUsage());
  }

  private String makeFlagUsage() {
    StringBuilder builder = new StringBuilder();
    for (CliFlag flag : shortFlagInfo.values()) {
      builder.append('[');
      builder.append(CliFlag.FLAG_START).append(flag.getShortFlag());
      if (flag.hasLongFlag()) {
        builder.append('|').append(CliFlag.FLAG_START).append(CliFlag.FLAG_START);
        builder.append(flag.getLongFlag());
      }
      if (flag.hasParameter()) {
        builder.append(' ').append('<').append(flag.getParameterName()).append('>');
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
    for (CliNodeImpl child : children.values()) {
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
    CliParseContext context = new CliParseContext();
    parse(args, 0, context);
    validateContext(context);
    runActiveNodeFromContext(context);
  }

  private int parse(String[] args, int index, CliParseContext context) {
    context.setActiveNode(this);
    while (index < args.length) {
      index = parseArg(args, index, context);
    }
    return index;
  }

  // Process the argument at index. Returns the next index to process.
  private int parseArg(String[] args, int index, CliParseContext context) {
    String arg = args[index];
    int nextIndex;
    if (isFlag(arg)) {
      nextIndex = parseFlag(args, index);
    } else if (isChild(arg) && context.getParameters().length == 0) {
      CliNodeImpl child = getChild(arg);
      nextIndex = child.parse(args, index + 1, context);
    } else {
      context.addParameter(arg);
      nextIndex = index + 1;
    }
    return nextIndex;
  }

  private boolean isFlag(String arg) {
    return (arg.charAt(0) == CliFlag.FLAG_START);
  }
  
  private int parseFlag(String[] args, int index) {
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
    boolean validFlag = (isLongFlag
                         ? longFlagShortFlags.containsKey(flagString)
                         : shortFlagInfo.containsKey(flagString));
    if (!validFlag) {
      throwBadArgException("Unknown flag '" + flagString + "'", args, index);
    }
    String shortFlag = (isLongFlag ? longFlagShortFlags.get(flagString) : flagString);
    CliFlag flag = shortFlagInfo.get(shortFlag);

    // Does it have an argument?
    List<String> arguments = new ArrayList<String>();
    if (flag.hasParameter()) {
      if (index == args.length - 1) {
        throwBadArgException("Expected argument for flag '" + flag + "'", args, index);
      }
      index += 1;
      arguments.add(args[index]);
    }

    flag.getAction().run(arguments.toArray(new String[0]));

    return index + 1;
  }

  private boolean isChild(String arg) {
    return children.containsKey(arg);
  }

  private CliNodeImpl getChild(String arg) {
    return children.get(arg);
  }

  private void throwBadArgException(String baseMessage, String[] args, int index) {
    String message = String.format("%s index=%d, arg=%s", baseMessage, index, args[index]);
    throw new IllegalArgumentException(message);
  }

  private void validateContext(CliParseContext context) {
    CliNodeImpl activeNode = context.getActiveNode();
    String[] parameters = context.getParameters();
    if (activeNode.isList() && parameters.length != 0) {
      throw new IllegalArgumentException("Unknown command '" + parameters[0]
                                         + "' for list " + activeNode.name);
    }
  }

  private void runActiveNodeFromContext(CliParseContext context) {
    CliNodeImpl activeNode = context.getActiveNode();
    String[] parameters = context.getParameters();
    activeNode.action.run(parameters);
  }
}
