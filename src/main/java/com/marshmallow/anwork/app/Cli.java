package com.marshmallow.anwork.app;

import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.Map;

/**
 * This is an object that can run actions based off of command line interface
 * arguments.
 *
 * This object parses strings of the following form.
 * <pre>
 *   -f (A short flag)
 *   -f something (A short flag and an argument)
 *   --flag (A long flag)
 *   --flag something (A long flag and an argument)
 * </pre>
 *
 * The class is meant to be used in the following way.
 * <pre>
 *   Cli cli = new Cli();
 *   cli.addAction("f", "flag", true, new CliAction() { ... });
 *   cli.addAction("v", null, false, new CliAction() { ... });
 *   cli.addAction("d", "debug", false, new CliAction() { ... });
 *   ...
 *   cli.run(args);
 * </pre>
 *
 * @author Andrew
 * @date Sep 9, 2017
 */
public class Cli {

  private static class FlagInfo {
    private final String shortFlag;
    private final String longFlag;
    private final String description;
    private final String argumentName;
    private final CliAction action;

    public FlagInfo(String shortFlag, String longFlag, String description, String argumentName, CliAction action) {
      this.shortFlag = shortFlag;
      this.longFlag = longFlag;
      this.description = description;
      this.argumentName = argumentName;
      this.action = action;
    }

    public String getShortFlag() {
      return shortFlag;
    }

    public boolean hasLongFlag() {
      return longFlag != null;
    }

    public String getLongFlag() {
      return longFlag;
    }

    public boolean hasDescription() {
      return description != null;
    }

    public String getDescription() {
      return description;
    }

    public boolean hasArgument() {
      return argumentName != null;
    }

    public String getArgumentName() {
      return argumentName;
    }

    public CliAction getAction() {
      return action;
    }
  }

  private static final char FLAG_START = '-';

  // This is a map from short flag to information about the flag.
  private Map<String, FlagInfo> shortFlagInfo = new LinkedHashMap<String, FlagInfo>();
  // This is a map from long flag to short flag.
  private Map<String, String> longFlagShortFlags = new HashMap<String, String>();

  /**
   * Add a command line action
   *
   * @param shortFlag The short flag, e.g., f, v, q, etc. This parameter is
   * required!
   * @param longFlag The long flag, e.g., flag, verbose, quiet, etc. This
   * parameter is optional. Pass <code>null</code> to ignore this parameter.
   * @param description The description for this command. This parameter is
   * optional. Pass <code>null</code> to ignore this parameter.
   * @param argumentName The name of the argument to this command. If <code>null
   * </code> is passed for this parameter, there will be no argument to this
   * command. This parameter is optional.
   * @param action The action to be run. This parameter is required!
   * @throws IllegalArgumentException when there is a required argument that is
   * passed as <code>null</code>
   * @see CliAction
   */
  public void addAction(String shortFlag,
                        String longFlag,
                        String description,
                        String argumentName,
                        CliAction action) throws IllegalArgumentException {
    if (shortFlag == null || action == null) {
      throw new IllegalArgumentException("Short flag and CLI action are required arguments!");
    }
    FlagInfo info = new FlagInfo(shortFlag, longFlag, description, argumentName, action);
    shortFlagInfo.put(shortFlag, info);
    if (longFlag != null) {
      longFlagShortFlags.put(longFlag, shortFlag);
    }
  }

  /**
   * Get the usage information for this CLI instance.
   *
   * @return The usage information for this CLI instance
   */
  public String getUsage() {
    StringBuilder builder = new StringBuilder("");
    for (FlagInfo flagInfo : shortFlagInfo.values()) {
      builder.append(FLAG_START).append(flagInfo.getShortFlag());
      if (flagInfo.hasLongFlag()) {
        builder.append('|').append(FLAG_START).append(FLAG_START).append(flagInfo.getLongFlag());
      }
      if (flagInfo.hasArgument()) {
        builder.append(' ').append('<').append(flagInfo.getArgumentName()).append('>');
      }
      if (flagInfo.hasDescription()) {
        builder.append(' ').append(flagInfo.getDescription());
      }
      builder.append('\n');
    }
    return builder.toString();
  }

  /**
   * Run the actions associated with the provided CLI arguments.
   *
   * @param args The command line arguments to be run
   * @throws IllegalArgumentException if there is incorrect CLI syntax
   */
  public void runActions(String[] args) throws IllegalArgumentException {
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
    if (arg.charAt(0) != FLAG_START) {
      throwBadArgException("Expected flag syntax", args, index);
    }
    if (arg.charAt(1) == FLAG_START) {
      isLongFlag = true;
    }

    // Is it a valid flag?
    String flag = arg.substring(isLongFlag ? 2 : 1);
    boolean validFlag = (isLongFlag ? longFlagShortFlags.containsKey(flag) : shortFlagInfo.containsKey(flag));
    if (!validFlag) {
      throwBadArgException("Unknown flag '" + flag + "'", args, index);
    }
    String shortFlag = (isLongFlag ? longFlagShortFlags.get(flag) : flag);
    FlagInfo flagInfo = shortFlagInfo.get(shortFlag);

    // Does it have an argument?
    String argument = null;
    if (flagInfo.hasArgument()) {
      if (index == args.length -1) {
        throwBadArgException("Expected argument for flag '" + flag + "'", args, index);
      }
      index += 1;
      argument = args[index];
    }

    flagInfo.getAction().run(argument);

    return index + 1;
  }

  private void throwBadArgException(String baseMessage, String args[], int index) throws IllegalArgumentException {
    String message = String.format("%s index=%d, arg=%s", baseMessage, index, args[index]);
    throw new IllegalArgumentException(message);
  }
}
