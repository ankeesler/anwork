package com.marshmallow.anwork.journal;

/**
 * This is an object that maintains a sequence of {@link JournalEntry}'s that have happened.
 *
 * <p>
 * Created Sep 17, 2017
 * </p>
 *
 * @author Andrew
 */
public interface Journal {

  /**
   * Add an entry to this journal.
   *
   * @param entry The entry to add to this journal
   */
  public void addEntry(JournalEntry entry);

  /**
   * Get the entries in this journal in the order in which they were added.
   *
   * @return The entries in this journal in the order in which they were added
   */
  public JournalEntry[] getEntries();
}
