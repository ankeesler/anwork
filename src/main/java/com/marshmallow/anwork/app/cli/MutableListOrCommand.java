package com.marshmallow.anwork.app.cli;

/**
 * This is a {@link ListOrCommand} that is editable, i.e., clients can set the information on a
 * {@link ListOrCommand} using this interface.
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
interface MutableListOrCommand extends ListOrCommand {

  /**
   * Set the name of this {@link ListOrCommand}.
   *
   * @param name The name set on this {@link ListOrCommand}
   * @return <code>this</code> MutableListOrCommand currently being edited
   */
  public MutableListOrCommand setName(String name);

  /**
   * Set the description of this {@link ListOrCommand}.
   *
   * @param description The description set on this {@link ListOrCommand}
   * @return <code>this</code> MutableListOrCommand currently being edited
   */
  public MutableListOrCommand setDescription(String description);


  /**
   * Add a {@link Flag} to this {@link ListOrCommand}.
   *
   * @param shortFlag The short flag to assign to the new {@link Flag}
   * @return A {@link Flag} that can edited, i.e., a {@link MutableFlag}
   */
  public MutableFlag addFlag(String shortFlag);
}