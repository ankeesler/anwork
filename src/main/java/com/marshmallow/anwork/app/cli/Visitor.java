package com.marshmallow.anwork.app.cli;

/**
 * This is an object that can visit the {@link Flag}, {@link List}, and {@link Command} members
 * of a {@link Cli}.
 *
 * <p>
 * Created Oct 1, 2017
 * </p>
 *
 * @author Andrew
 */
public interface Visitor {

  /**
   * Visit a {@link Cli} {@link Flag}.
   *
   * @param flag The {@link Flag} that is currently being visited
   */
  public void visitFlag(Flag flag);

  /**
   * Visit a {@link Cli} {@link List}.
   *
   * @param list The {@link List} currently being visited
   */
  public void visitList(List list);

  /**
   * Leave a {@link Cli} {@link List}.
   *
   * <p>
   * This is a helpful utility function to determine where a {@link List} starts and ends in a
   * {@link Cli}.
   * </p>
   *
   * @param list The {@link Cli} {@link List} that is being left
   */
  public void leaveList(List list);

  /**
   * Visit a {@link Cli} {@link Flag}.
   *
   * @param command The {@link Command} currently being visited
   */
  public void visitCommand(Command command);
}
