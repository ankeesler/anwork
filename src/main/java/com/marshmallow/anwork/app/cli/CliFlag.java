package com.marshmallow.anwork.app.cli;

/**
 * This is a single flag passed to a command line invocation.
 *
 * This object comes in a couple of different formats.
 * <pre>
 *   -f (A short flag) (see {@link #makeShortFlag(String, String, CliAction)})
 *   -f something (A short flag and an argument) (see {@link #makeShortFlagWithParameter(String, String, String, CliAction)}
 *   --flag (A long flag) (see {@link #makeLongFlag(String, String, String, CliAction)})
 *   --flag something (A long flag and an argument) (see {@link #makeLongFlagWithParameter(String, String, String, String, CliAction)})
 * </pre>
 *
 * @author Andrew
 * @date Sep 9, 2017
 */
public class CliFlag {

  public static final char FLAG_START = '-';

  /**
   * Make a short flag.
   *
   * @param shortFlag The name of the short flag, i.e., "d", "v", "o", etc.
   * @param description The description of the flag
   * @param action The action that the flag takes
   * @return A new instance of a CliFlag
   */
  public static CliFlag makeShortFlag(String shortFlag,
                                      String description,
                                      CliAction action) {
    return new CliFlag(shortFlag, null, description, null, action);
  }

  /**
   * Make a short flag that takes a parameter.
   *
   * @param shortFlag The name of the short flag, i.e., "d", "v", "o", etc.
   * @param description The description of the flag
   * @param parameterName The name of the parameter that the flag takes 
   * @param action The action that the flag takes
   * @return A new instance of a CliFlag
   */
  public static CliFlag makeShortFlagWithParameter(String shortFlag,
                                                   String description,
                                                   String parameterName,
                                                   CliAction action) {
    return new CliFlag(shortFlag, null, description, parameterName, action);
  }

  /**
   * Make a long flag.
   *
   * @param shortFlag The name of the short flag, i.e., "d", "v", "o", etc.
   * @param longFlag The name of the long flag, i.e., "debug", "verbose",
   * "output", etc.
   * @param description The description of the flag
   * @param action The action that the flag takes
   * @return A new instance of a CliFlag
   */
  public static CliFlag makeLongFlag(String shortFlag,
                                     String longFlag,
                                     String description,
                                     CliAction action) {
    return new CliFlag(shortFlag, longFlag, description, null, action);
  }

  /**
   * Make a long flag that takes a parameter.
   *
   * @param shortFlag The name of the short flag, i.e., "d", "v", "o", etc.
   * @param longFlag The name of the long flag, i.e., "debug", "verbose",
   * "output", etc.
   * @param description The description of the flag
   * @param parameterName The name of the parameter that the flag takes 
   * @param action The action that the flag takes
   * @return A new instance of a CliFlag
   */
  public static CliFlag makeLongFlagWithParameter(String shortFlag,
                                                  String longFlag,
                                                  String description,
                                                  String parameterName,
                                                  CliAction action) {
    return new CliFlag(shortFlag, longFlag, description, parameterName, action);
  }

  private final String shortFlag;
  private final String longFlag;
  private final String description;
  private final String parameterName;
  private final CliAction action;

  // See static "make" methods above.
  private CliFlag(String shortFlag, String longFlag, String description, String parameterName, CliAction action) {
    this.shortFlag = shortFlag;
    this.longFlag = longFlag;
    this.description = description;
    this.parameterName = parameterName;
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

  public String getDescription() {
    return description;
  }

  public boolean hasParameter() {
    return parameterName != null;
  }

  public String getParameterName() {
    return parameterName;
  }

  public CliAction getAction() {
    return action;
  }
}
