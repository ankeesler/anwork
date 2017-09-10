package com.marshmallow.anwork.app.cli;

/**
 * This is an action that can be run via a CLI command.
 *
 * @author Andrew
 * Created Sep 9, 2017
 */
public interface CliAction {

  /**
   * Run the CLI action.
   *
   * @param parameters The parameters to the action; this parameter is never
   * <code>null</code>
   * @see CliNode#addShortFlag(String, String, CliAction)
   * @see CliNode#addShortFlagWithParameter(String, String, String, CliAction)
   * @see CliNode#addLongFlag(String, String, String, CliAction)
   * @see CliNode#addLongFlagWithParameter(String, String, String, String, CliAction)
   * @see CliNode#addCommand(String, String, CliAction)
   */
  public void run(String[] parameters);

}
