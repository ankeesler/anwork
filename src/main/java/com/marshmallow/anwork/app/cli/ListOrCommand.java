package com.marshmallow.anwork.app.cli;

/**
 * This is a base interface that represents both a {@link List} and a {@link Command}.
 *
 * <p>
 * Both {@link List}'s and {@link Command}'s are required to have a name and can optionally
 * have a description. Furthermore, both {@link List}'s and {@link Command}'s can have
 * {@link Flag}'s.
 * </p>
 *
 * <p>
 * Note: this interface is package-private on purpose. Public clients are meant to use the
 * {@link List} and {@link Command} interfaces.
 * </p>
 *
 * <p>
 * Created Oct 14, 2017
 * </p>
 *
 * @author Andrew
 */
interface ListOrCommand {

  /**
   * Get the name of this {@link ListOrCommand}.
   *
   * @return The name of this {@link ListOrCommand}
   */
  public String getName();

  /**
   * Get whether or not this {@link ListOrCommand} has a description.
   *
   * @return Whether or not this {@link ListOrCommand} has a description
   */
  public boolean hasDescription();

  /**
   * Get the description for this {@link ListOrCommand}, if it exists. If no description exists,
   * the returned data is undefined. Use {@link #hasDescription()} to tell whether or not the
   * description exists.
   *
   * @return The description for this {@link ListOrCommand}, if it exists; if no description
   *     exists, the returned data is undefined
   */
  public String getDescription();

  /**
   * Get the {@link Flag}'s on this {@link ListOrCommand}.
   *
   * @return The {@link Flag}'s on this {@link ListOrCommand}
   */
  public Flag[] getFlags();
}
