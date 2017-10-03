package com.marshmallow.anwork.app.cli;

/**
 * This is the base interface for how to work with a node in a CLI tree. 
 *
 * <p>
 * This interface is to be used within this package. We may expose this in the future, if we
 * want. It was originally created at package-scope to simplify the public interface to the CLI
 * system.
 * </p>
 *
 * <p>
 * Created Sep 10, 2017
 * </p>
 *
 * @author Andrew
 * @see CliList
 * @see CliCommand
 */
interface CliNode {

  /**
   * Make a short flag.
   *
   * <p>
   * Note: this flag is created with the {@link CliArgumentType} of
   * {@link CliArgumentType#BOOLEAN}.
   * </p>
   *
   * @param shortFlag The name of the short flag, i.e., "d", "v", "o", etc.
   * @param description The description of the flag
   */
  public void addShortFlag(String shortFlag, String description);

  /**
   * Make a short flag that takes a parameter.
   *
   * @param shortFlag The name of the short flag, i.e., "d", "v", "o", etc.
   * @param description The description of the flag
   * @param parameterName The name of the parameter that the flag takes, e.g., "file path,"
   *     "directory," ...
   * @param parameterName The description of the parameter that the flag takes
   * @param parameterType The {@link CliArgumentType} of the parameter
   */
  public void addShortFlagWithParameter(String shortFlag,
                                        String description,
                                        String parameterName,
                                        String parameterDescription,
                                        CliArgumentType parameterType);

  /**
   * Make a long flag.
   *
   * <p>
   * Note: this flag is created with the {@link CliArgumentType} of
   * {@link CliArgumentType#BOOLEAN}.
   * </p>
   *
   * @param shortFlag The name of the short flag, i.e., "d", "v", "o", etc.
   * @param longFlag The name of the long flag, i.e., "debug", "verbose",
   *     "output", etc.
   * @param description The description of the flag
   */
  public void addLongFlag(String shortFlag,
                          String longFlag,
                          String description);

  /**
   * Make a long flag that takes a parameter.
   *
   * @param shortFlag The name of the short flag, i.e., "d", "v", "o", etc.
   * @param longFlag The name of the long flag, i.e., "debug", "verbose",
   *     "output", etc.
   * @param description The description of the flag
   * @param parameterName The name of the parameter that the flag takes, e.g., "file path,"
   *     "directory," ...
   * @param parameterType The {@link CliArgumentType} of the parameter
   */
  public void addLongFlagWithParameter(String shortFlag,
                                       String longFlag,
                                       String description,
                                       String parameterName,
                                       String parameterDescription,
                                       CliArgumentType parameterType);
}
