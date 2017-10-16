package com.marshmallow.anwork.app.cli;

/**
 * This is a single flag passed to a command line invocation.
 *
 * <p>A {@link Flag} can represent a couple different types of command line items.
 * <pre>
 *   -f (A short flag)
 *   -f something (A short flag and an argument)
 *   --flag (A long flag)
 *   --flag something (A long flag and an argument)
 * </pre>
 *
 * <p>A {@link Flag} is required to contain the following information.
 * <ul>
 * <li>A short flag, e.g., "d", "v", "o", etc.</li>
 * </ul>
 * A {@link Flag} can optionally contain the following information.
 * <ul>
 * <li>A long flag, e.g., "debug", "verbose", "output", etc.</li>
 * <li>A description of what the flag means</li>
 * <li>A {@link Argument} passed to the flag, e.g., "--output path/to/file", "-f tuna.txt"</li>
 * </ul>
 *
 * <p>
 * Created Sep 9, 2017
 * </p>
 *
 * @author Andrew
 */
public interface Flag extends Comparable<Flag> {

  /**
   * This is the prefix to a flag passed on the command line. All flags must start with this
   * prefix.
   */
  public static final char FLAG_START = '-';

  /**
   * Get this {@link Flag}'s short flag, e.g., "d", "v", "o", etc.
   *
   * @return This {@link Flag}'s short flag
   */
  public String getShortFlag();

  /**
   * Get whether or not this {@link Flag} has a long flag associated with it.
   *
   * @return Whether or not this {@link Flag} has a long flag associated with it
   */
  public boolean hasLongFlag();

  /**
   * Get this {@link Flag}'s long flag, if it exists. If no long flag exists, the returned data is
   * undefined. Use {@link #hasLongFlag()} to determine whether or not a long flag exists.
   *
   * @return This {@link Flag}'s long flag, if it exists; if no long flag exists, the returned data
   *     is undefined
   */
  public String getLongFlag();

  /**
   * Get whether or not this {@link Flag} has a description.
   *
   * @return Whether or not this {@link Flag} has a description
   */
  public boolean hasDescription();

  /**
   * Get the description for this {@link Flag}, if it exists. If no description exists, the
   * returned data is undefined. Use {@link #hasDescription()} to tell whether or not the
   * description exists.
   *
   * @return The description for this {@link Flag}, if it exists; if no description exists, the
   *     returned data is undefined
   */
  public String getDescription();

  /**
   * Get whether or not this {@link Flag} has an {@link Argument} passed to it.
   *
   * @return Whether or not this {@link Flag} has an {@link Argument} passed to it
   */
  public boolean hasArgument();

  /**
   * Get the {@link Argument} for this {@link Flag}, if it exists. If no {@link Argument} exists,
   * the returned data is undefined. Use {@link #hasDescription()} to tell whether or not the
   * {@link Argument} exists.
   *
   * @return The {@link Argument} for this {@link Flag}, if it exists; if no {@link Argument}
   *     exists, the returned data is undefined
   */
  public Argument getArgument();
}
