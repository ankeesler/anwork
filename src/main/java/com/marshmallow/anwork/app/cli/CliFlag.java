package com.marshmallow.anwork.app.cli;

/**
 * This is a single flag passed to a command line invocation.
 *
 * This object can represent a couple different types of flags.
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
class CliFlag {

  public static final char FLAG_START = '-';

  static CliFlag makeShortFlag(String shortFlag,
                               String description,
                               CliAction action) {
    return new CliFlag(shortFlag, null, description, null, action);
  }

  static CliFlag makeShortFlagWithParameter(String shortFlag,
                                            String description,
                                            String parameterName,
                                            CliAction action) {
    return new CliFlag(shortFlag, null, description, parameterName, action);
  }

  static CliFlag makeLongFlag(String shortFlag,
                              String longFlag,
                              String description,
                              CliAction action) {
    return new CliFlag(shortFlag, longFlag, description, null, action);
  }

  static CliFlag makeLongFlagWithParameter(String shortFlag,
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

  String getShortFlag() {
    return shortFlag;
  }

  boolean hasLongFlag() {
    return longFlag != null;
  }

  String getLongFlag() {
    return longFlag;
  }

  String getDescription() {
    return description;
  }

  boolean hasParameter() {
    return parameterName != null;
  }

  String getParameterName() {
    return parameterName;
  }

  CliAction getAction() {
    return action;
  }
}
