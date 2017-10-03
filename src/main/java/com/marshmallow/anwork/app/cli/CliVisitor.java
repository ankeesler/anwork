package com.marshmallow.anwork.app.cli;

/**
 * This is an object that can visit the members of a {@link Cli} tree.
 *
 * <p>
 * Created Oct 1, 2017
 * </p>
 *
 * @author Andrew
 */
public interface CliVisitor {

  /**
   * Visit a short flag.
   *
   * @param shortFlag The name of the short flag, i.e., "d", "v", "o", etc.
   * @param description The description of the flag
   */
  public void visitShortFlag(String shortFlag,
                             String description);

  /**
   * Visit a short flag that has a parameter.
   *
   * @param shortFlag The name of the short flag, i.e., "d", "v", "o", etc.
   * @param description The description of the flag
   * @param parameterName The name of the parameter that the flag takes, e.g., "file path,"
   *     "directory," ...
   * @param parameterDescription The description of the parameter
   * @param parameterType The {@link CliArgumentType} of the parameter
   */
  public void visitShortFlagWithParameter(String shortFlag,
                                          String description,
                                          String parameterName,
                                          String parameterDescription,
                                          CliArgumentType parameterType);

  /**
   * Visit a long flag.
   *
   * @param shortFlag The name of the short flag, i.e., "d", "v", "o", etc.
   * @param longFlag The name of the long flag, i.e., "debug", "verbose",
   *     "output", etc.
   * @param description The description of the flag
   */
  public void visitLongFlag(String shortFlag,
                            String longFlag,
                            String description);

  /**
   * Visit a long flag that has a parameter.
   *
   * @param shortFlag The name of the short flag, i.e., "d", "v", "o", etc.
   * @param longFlag The name of the long flag, i.e., "debug", "verbose",
   *     "output", etc.
   * @param description The description of the flag
   * @param parameterName The name of the parameter that the flag takes, e.g., "file path,"
   *     "directory," ...
   * @param parameterDescription The description of the parameter
   * @param parameterType The {@link CliArgumentType} of the parameter
   */
  public void visitLongFlagWithParameter(String shortFlag,
                                         String longFlag,
                                         String description,
                                         String parameterName,
                                         String parameterDescription,
                                         CliArgumentType parameterType);

  /**
   * Visit a CLI list.
   *
   * @param name The name of the list
   * @param description The description of the list
   */
  public void visitList(String name, String description);

  /**
   * Leave a CLI list.
   *
   * <p>
   * This is a helpful utility function to determine where a list starts and ends.
   * </p>
   *
   * @param name The name of the list
   */
  public void leaveList(String name);

  /**
   * Visit a CLI command.
   *
   * @param name The name of the command
   * @param description The description of the command
   */
  public void visitCommand(String name, String description);
}
