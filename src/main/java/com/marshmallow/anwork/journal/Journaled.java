package com.marshmallow.anwork.journal;

/**
 * This is a type of object that has a {@link Journal} associated with it.
 *
 * <p>
 * Created Sep 17, 2017
 * </p>
 *
 * @author Andrew
 */
public interface Journaled {

  /**
   * Get the journal associated with this object.
   *
   * @return The journal associated with this object
   */
  public Journal getJournal();
}
