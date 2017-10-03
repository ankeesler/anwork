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
   * @param parameters The parameters to the action; this parameter is never <code>null</code>
   * @param flags The {@link CliFlags} passed to this CLI command
   * @see CliNodeImpl#addCommand(String, String, CliAction)
   */
  public void run(CliFlags flags, String[] parameters);

}
