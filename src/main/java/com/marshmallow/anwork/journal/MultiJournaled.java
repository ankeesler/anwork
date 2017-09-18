package com.marshmallow.anwork.journal;

/**
 * This is an object that has more than one associated {@link Journal}'s. These {@link Journal}'s
 * are accessed via a key.
 *
 * <p>
 * Created Sep 17, 2017
 * </p>
 *
 * @author Andrew
 */
public interface MultiJournaled extends Journaled {

  /**
   * Get the journal associated with this object using the key provided.
   *
   * @param key The key to use to fetch a specific journal associated with this object
   * @return The journal associated with this object using the key provided
   */
  public Journal getJournal(String key);
}
