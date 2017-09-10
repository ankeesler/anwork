package com.marshmallow.anwork.app.cli;

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
   * @param parameters The parameters to the action; this parameter is never
   * <code>null</code>
   * @see CliFlag#makeShortFlag(String, String, String, CliAction)
   * @see CliFlag#makeShortFlagWithParameter(String, String, String, CliAction)
   * @see CliFlag#makeLongFlag(String, String, String, CliAction)
   * @see CliFlag#makeLongFlagWithParameter(String, String, String, String, CliAction)
   * @see CliNode#makeCommand(String, String, CliAction)
   */
  public void run(String[] parameters);

}
