package com.marshmallow.anwork.app.cli;

/**
 * This is a {@link List} that is editable, i.e., clients can set the information on a
 * {@link List} using this interface.
 *
 * <p>
 * Created Oct 14, 2017
 * </p>
 *
 * @author Andrew
 */
public interface MutableList extends MutableListOrCommand, List {

  @Override
  public MutableList setName(String name);

  @Override
  public MutableList setDescription(String description);

  /**
   * Add a {@link Command} to this {@link List}.
   *
   * @param name The name of the {@link Command}
   * @param action The {@link Action} that the {@link Command} should take
   * @return An editable {@link Command}, i.e., an {@link MutableCommand}
   */
  public MutableCommand addCommand(String name, Action action);

  /**
   * Add a {@link List} to this {@link List}.
   *
   * @param name The name of the {@link List}
   * @return An editable {@link List}, i.e., an {@link MutableList}
   */
  public MutableList addList(String name);
}
