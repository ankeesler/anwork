package com.marshmallow.anwork.app.cli;

/**
 * This is a single flag passed to a command line invocation.
 *
 * <p>
 * This object can represent a couple different types of flags.
 * <pre>
 *   -f (A short flag)
 *   (see {@link #makeShortFlag(String, String, CliAction)})
 * 
 *   -f something (A short flag and an argument)
 *   (see {@link #makeShortFlagWithParameter(String, String, String, CliAction)}
 *
 *   --flag (A long flag)
 *   (see {@link #makeLongFlag(String, String, String, CliAction)})
 * 
 *   --flag something (A long flag and an argument)
 *   (see {@link #makeLongFlagWithParameter(String, String, String, String, CliAction)})
 * </pre>
 * </p>
 *
 * <p>
 * Created Sep 9, 2017
 * </p>
 *
 * @author Andrew
 */
class CliFlag implements Comparable<CliFlag> {

  public static final char FLAG_START = '-';

  static CliFlag makeShortFlag(String shortFlag,
                               String description) {
    return new CliFlag(shortFlag, null, description, null, null, null);
  }

  static CliFlag makeShortFlagWithParameter(String shortFlag,
                                            String description,
                                            String parameterName,
                                            String parameterDescription,
                                            CliArgumentType parameterType) {
    return new CliFlag(shortFlag,
                       null, // longFlag
                       description,
                       parameterName,
                       parameterDescription,
                       parameterType);
  }

  static CliFlag makeLongFlag(String shortFlag,
                              String longFlag,
                              String description) {
    return new CliFlag(shortFlag, longFlag, description, null, null, null);
  }

  static CliFlag makeLongFlagWithParameter(String shortFlag,
                                           String longFlag,
                                           String description,
                                           String parameterName,
                                           String parameterDescription,
                                           CliArgumentType parameterType) {
    return new CliFlag(shortFlag,
                       longFlag,
                       description,
                       parameterName,
                       parameterDescription,
                       parameterType);
  }

  private final String shortFlag;
  private final String longFlag;
  private final String description;
  private final String parameterName;
  private final String parameterDescription;
  private final CliArgumentType parameterType;

  // See static "make" methods above.
  private CliFlag(String shortFlag,
                  String longFlag,
                  String description,
                  String parameterName,
                  String parameterDescription,
                  CliArgumentType parameterType) {
    this.shortFlag = shortFlag;
    this.longFlag = longFlag;
    this.description = description;
    this.parameterName = parameterName;
    this.parameterDescription = parameterDescription;
    this.parameterType = parameterType;
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

  String getParameterDescription() {
    return parameterDescription;
  }

  CliArgumentType getParameterType() {
    return parameterType;
  }

  @Override
  public int compareTo(CliFlag other) {
    return shortFlag.compareTo(other.shortFlag);
  }
}