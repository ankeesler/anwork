package com.marshmallow.anwork.app.cli;

/**
 * This is an {@link Argument} that is editable, i.e., clients can set the information on a
 * {@link Argument} using this interface.
 *
 * <p>
 * Created Oct 15, 2017
 * </p>
 *
 * @author Andrew
 */
public interface MutableArgument extends Argument {

  /**
   * Set the name of this {@link Argument}.
   *
   * @param name The name to set on this {@link Argument}
   * @return <code>this</code> current {@link MutableArgument} currently being edited
   */
  public MutableArgument setName(String name);

  /**
   * Set the {@link ArgumentType} of this {@link Argument}.
   *
   * @param <T> The backing Java type for the passed {@link ArgumentType}
   * @param type The {@link ArgumentType} to set on this {@link Argument}.
   * @return <code>this</code> current {@link MutableArgument} currently being edited
   */
  public <T> MutableArgument setType(ArgumentType<T> type);

  /**
   * Set the description of this {@link Argument}.
   *
   * @param description The description to set on this {@link Argument}.
   * @return <code>this</code> current {@link MutableArgument} currently being edited
   */
  public MutableArgument setDescription(String description);
}
