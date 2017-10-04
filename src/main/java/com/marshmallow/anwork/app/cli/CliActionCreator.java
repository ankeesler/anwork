package com.marshmallow.anwork.app.cli;

/**
 * An instance of this class simply creates a {@link CliAction} given a {@link CliCommand} name.
 *
 * <p>
 * Created Oct 4, 2017
 * </p>
 *
 * @author Andrew
 */
public interface CliActionCreator {

  /**
   * Create a {@link CliAction} given a {@link CliCommand} name.
   *
   * @param commandName The name of the {@link CliCommand} for which to create the
   *     {@link CliAction}
   * @return A {@link CliAction} for the provided {@link CliCommand} name
   */
  public CliAction createAction(String commandName);
}