package com.marshmallow.anwork.app.cli;

/**
 * This is a {@link Flag} that is editable, i.e., clients can set the information on a
 * {@link Flag} using this interface.
 *
 * <p>
 * Created Oct 14, 2017
 * </p>
 *
 * @author Andrew
 */
public interface MutableCommand extends MutableListOrCommand, Command {

  @Override
  public MutableCommand setName(String name);

  @Override
  public MutableCommand setDescription(String description);

  /**
   * Set this {@link Command}'s {@link Action}.
   *
   * @param action The {@link Action} to be set on this {@link Command}
   * @return <code>this</code> {@link Command} currently being edited
   */
  public MutableCommand setAction(Action action);

  /**
   * Add an {@link Argument} to this {@link Command}. Note that the order in which these
   * {@link Argument}'s are added matters. If a {@link Command} "tuna" has 3 {@link Argument}'s
   * of type {@link ArgumentType#STRING}, {@link ArgumentType#NUMBER}, and
   * {@link ArgumentType#STRING}, then they will be validated that way upon
   * {@link Cli#parse(String[])}.
   *
   * @param <T> the backing Java type of the {@link ArgumentType} for this new {@link Argument}
   * @param name The name of the new {@link Argument} to add to this {@link Command}
   * @param type The {@link ArgumentType} of the new {@link Argument} to add to this
   *     {@link Command}
   * @return The {@link Argument} created from calling this method
   */
  public <T> MutableArgument addArgument(String name, ArgumentType<T> type);
}
