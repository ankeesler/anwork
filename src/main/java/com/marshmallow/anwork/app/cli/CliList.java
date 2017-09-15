package com.marshmallow.anwork.app.cli;

/**
 * This is a specific node in the CLI tree that is a list of other lists
 * and commands.
 *
 * <p>
 * This interface is to be used publicly.
 * </p>
 *
 * <p>
 * Created Sep 10, 2017
 * </p>
 *
 * @author Andrew
 * @see CliCommand
 */
public interface CliList extends CliNode {

  /**
   * Add a new command list.
   *
   * <p>
   * A command list may be something like "task" with child commands "create,"
   * "update," and "delete."
   * </p>
   *
   * @param name The name of the command list
   * @param description The description of the command list
   * @return The new CLI node that represents the command list
   * @see #addCommand(String, String, CliAction)
   */
  public CliList addList(String name, String description);

  /**
   * Add a new command.
   *
   * <p>
   * A command may be something like "create," "update," or "delete."
   * </p>
   *
   * @param name The name of the command
   * @param description The description of the command
   * @param action The action associated with the command
   * @return The new CLI node representing the command
   * @see #addList(String, String)
   */
  public CliCommand addCommand(String name, String description, CliAction action);

}
