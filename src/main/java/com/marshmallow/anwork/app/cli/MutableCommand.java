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
}
