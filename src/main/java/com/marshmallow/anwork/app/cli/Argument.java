package com.marshmallow.anwork.app.cli;

/**
 * This is an argument passed to a command line {@link Flag} or {@link Command}. An instance of
 * this type is required to have a name and a type, and can optionally have a description.
 *
 * <p>
 * Created Oct 14, 2017
 * </p>
 *
 * @author Andrew
 * @see ArgumentType
 */
public interface Argument {

  /**
   * Get the name for this {@link Argument}, e.g., "file-path", "output-directory", "count", etc.
   *
   * @return The name for this {@link Argument}
   */
  public String getName();

  /**
   * Get the type of this {@link Argument}.
   *
   * @return The type of this {@link Argument}
   */
  public ArgumentType<?> getType();

  /**
   * Get whether or not this {@link Argument} has a description.
   *
   * @return Whether or not this {@link Argument} has a description
   */
  public boolean hasDescription();

  /**
   * Get the description for this {@link Argument}, if it exists. If no description exists, the
   * returned data is undefined. Use {@link #hasDescription()} to tell whether or not the
   * description exists.
   *
   * @return The description for this {@link Argument}, if it exists; if no description exists, the
   *     returned data is undefined
   */
  public String getDescription();
}
