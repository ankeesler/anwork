package com.marshmallow.anwork.app.cli;

/**
 * This is a {@link Flag} which allows clients to edit its information.
 *
 * <p>
 * Created Oct 14, 2017
 * </p>
 *
 * @author Andrew
 */
public interface MutableFlag extends Flag {

  /**
   * Set this {@link Flag}'s short flag.
   *
   * @param shortFlag The short flag to use for this {@link Flag}
   * @return <code>this</code> {@link MutableFlag} currently being edited
   */
  public MutableFlag setShortFlag(String shortFlag);

  /**
   * Set this {@link Flag}'s longflag.
   *
   * @param longFlag The long flag to use for this {@link Flag}
   * @return <code>this</code> {@link MutableFlag} currently being edited
   */
  public MutableFlag setLongFlag(String longFlag);

  /**
   * Set this {@link Flag}'s description.
   *
   * @param description The description to use for this {@link Flag}
   * @return <code>this</code> {@link MutableFlag} currently being edited
   */
  public MutableFlag setDescription(String description);

  /**
   * Set this {@link Flag}'s {@link Argument}.
   *
   * @param <T> The backing Java type for the passed {@link ArgumentType}
   * @param name The name of the {@link Argument} to use for this {@link Flag}
   * @param type The type of the {@link Argument} to use for this {@link Flag}
   * @return <code>this</code> new {@link MutableArgument} currently being added
   */
  public <T> MutableArgument setArgument(String name, ArgumentType<T> type);
}
