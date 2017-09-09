package com.marshmallow.anwork.app.cli;

/**
 * This is a node in a CLI tree.
 *
 * A node has some flags and some children. It can run actions based off of
 * command line interface arguments (see {@link #parse(String[])}). 
 *
 * @author Andrew
 * @date Sep 9, 2017
 */
public interface CliNode {

  /**
   * Add a flag to this node.
   *
   * @param flag The flag to add to this node.
   */
  public void addFlag(CliFlag flag);

  /**
   * Add a command to this node.
   *
   * @param name The name of the command
   * @param description The description of the command
   * @param action The action for the commmand
   * @return The CLI node for the new command
   */
  public CliNode addCommand(String name, String description, CliAction action);

  /**
   * Parse an array of arguments.
   *
   * @param args An array of arguments
   */
  public void parse(String[] args);

  /**
   * Get the name of this CLI node.
   *
   * @return The name of this CLI node
   */
  public String getName();

  /**
   * Get the description of this CLI node.
   *
   * @return The description of this CLI node
   */
  public String getDescription();

  /**
   * Get a string describing this CLI node.
   *
   * @return A string describing this CLI node
   */
  public String getUsage();
}
