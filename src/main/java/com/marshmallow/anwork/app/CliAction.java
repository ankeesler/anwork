package com.marshmallow.anwork.app;

/**
 * This is an action that can be run via a CLI command.
 *
 * @author Andrew
 * @date Sep 9, 2017
 */
public interface CliAction {

  /**
   * Run the CLI action.
   *
   * @param argument The argument to the action, or <code>null</code> if there
   * is no argument to the action.
   * @see Cli#addAction(String, String, boolean, CliAction)
   */
  public void run(String argument);

}
