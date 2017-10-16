package com.marshmallow.anwork.app.cli;

/**
 * This is a collection of {@link Command}'s in a CLI object. The collection is identified by a
 * name (see {@link #getName()}). The collection can have {@link Flag}'s, {@link Command}'s and
 * other {@link List}'s added to it. The collection can optionally have a description.
 *
 * <p>
 * Created Oct 14, 2017
 * </p>
 *
 * @author Andrew
 */
public interface List extends ListOrCommand {

  /**
   * Get the {@link Command}'s in this {@link List}.
   *
   * @return The {@link Command}'s in this {@link List}, in no particular order
   */
  public Command[] getCommands();

  /**
   * Get the {@link List}'s in this {@link List}.
   *
   * @return The {@link List}'s in this {@link List}, in no particular order
   */
  public List[] getLists();
}
