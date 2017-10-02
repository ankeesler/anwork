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
   * @param shortFlag The name of the short flag
   * @param description The description of the short flag
   */
  public void visitShortFlag(String shortFlag,
                             String description);

  /**
   * Visit a short flag that has a parameter.
   *
   * @param shortFlag The name of the short flag
   * @param parameterName The name of the parameter
   * @param description The description of the short flag
   */
  public void visitShortFlagWithParameter(String shortFlag,
                                          String parameterName,
                                          String description);

  /**
   * Visit a long flag.
   *
   * @param shortFlag The name of the short flag
   * @param longFlag The name of the long flag
   * @param description The description of the long flag
   */
  public void visitLongFlag(String shortFlag,
                            String longFlag,
                            String description);

  /**
   * Visit a long flag that has a parameter.
   *
   * @param shortFlag The name of the short flag
   * @param longFlag The name of the long flag
   * @param parameterName The name of the parameter
   * @param description The description of the long flag
   */
  public void visitLongFlagWithParameter(String shortFlag,
                                         String longFlag,
                                         String parameterName,
                                         String description);

  /**
   * Visit a CLI list.
   *
   * @param name The name of the list
   * @param description The description of the list
   */
  public void visitList(String name, String description);

  /**
   * Visit a CLI command.
   *
   * @param name The name of the command
   * @param description The description of the command
   */
  public void visitCommand(String name, String description);
}
