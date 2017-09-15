package com.marshmallow.anwork.app.cli;

/**
 * This is an action that can be run via a CLI command.
 *
 * <p>
 * Created Sep 9, 2017
 * </p>
 *
 * @author Andrew
 */
public interface CliAction {

  /**
   * Run the CLI action.
   *
   * @param parameters The parameters to the action; this parameter is never
   * <code>null</code>
   * @see CliNodeImpl#addShortFlag(String, String, CliAction)
   * @see CliNodeImpl#addShortFlagWithParameter(String, String, String, CliAction)
   * @see CliNodeImpl#addLongFlag(String, String, String, CliAction)
   * @see CliNodeImpl#addLongFlagWithParameter(String, String, String, String, CliAction)
   * @see CliNodeImpl#addCommand(String, String, CliAction)
   */
  public void run(String[] parameters);

}
