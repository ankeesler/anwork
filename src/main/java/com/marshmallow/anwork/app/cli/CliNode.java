package com.marshmallow.anwork.app.cli;

/**
 * This is the base interface for how to work with a node in a CLI tree. 
 *
 * <p>
 * This interface is to be used within this package.
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
   * @param shortFlag The name of the short flag, i.e., "d", "v", "o", etc.
   * @param description The description of the flag
   * @param action The action that the flag takes
   */
  public void addShortFlag(String shortFlag, String description, CliAction action);

  /**
   * Make a short flag that takes a parameter.
   *
   * @param shortFlag The name of the short flag, i.e., "d", "v", "o", etc.
   * @param description The description of the flag
   * @param parameterName The name of the parameter that the flag takes
   * @param action The action that the flag takes
   */
  public void addShortFlagWithParameter(String shortFlag,
                                        String description,
                                        String parameterName,
                                        CliAction action);

  /**
   * Make a long flag.
   *
   * @param shortFlag The name of the short flag, i.e., "d", "v", "o", etc.
   * @param longFlag The name of the long flag, i.e., "debug", "verbose",
   *     "output", etc.
   * @param description The description of the flag
   * @param action The action that the flag takes
   */
  public void addLongFlag(String shortFlag,
                          String longFlag,
                          String description,
                          CliAction action);

  /**
   * Make a long flag that takes a parameter.
   *
   * @param shortFlag The name of the short flag, i.e., "d", "v", "o", etc.
   * @param longFlag The name of the long flag, i.e., "debug", "verbose",
   *     "output", etc.
   * @param description The description of the flag
   * @param parameterName The name of the parameter that the flag takes
   * @param action The action that the flag takes
   */
  public void addLongFlagWithParameter(String shortFlag,
                                       String longFlag,
                                       String description,
                                       String parameterName,
                                       CliAction action);
}
