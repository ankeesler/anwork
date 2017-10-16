package com.marshmallow.anwork.app.cli;

/**
 * An instance of this class simply creates a {@link Action} given a {@link Command} name.
 *
 * <p>
 * Created Oct 4, 2017
 * </p>
 *
 * @author Andrew
 */
public interface ActionCreator {

  /**
   * Create a {@link Action} given a {@link Command} name.
   *
   * @param commandName The name of the {@link Command} for which to create the
   *     {@link Action}
   * @return A {@link Action} for the provided {@link Command} name
   */
  public Action createAction(String commandName);
}