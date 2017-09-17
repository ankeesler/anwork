package com.marshmallow.anwork.journal;

import java.util.function.Predicate;
import java.util.stream.Stream;

/**
 * This is a {@link Journal} whose entries depend on a predicate applied to another
 * {@link Journal}'s entries.
 *
 * <p>
 * Created Sep 17, 2017
 * </p>
 *
 * @author Andrew
 */
public class FilteredJournal implements Journal {

  private final Journal journal;
  private final Predicate<JournalEntry> predicate;

  /**
   * Create a filtered journal given an underlying journal and a predicate.
   *
   * @param journal The underlying journal to use
   * @param predicate The predicate to apply to the entries in the underlying journal
   */
  public FilteredJournal(Journal journal, Predicate<JournalEntry> predicate) {
    this.journal = journal;
    this.predicate = predicate;
  }

  @Override
  public void addEntry(JournalEntry entry) {
    journal.addEntry(entry);
  }

  @Override
  public JournalEntry[] getEntries() {
    return Stream.of(journal.getEntries()).filter(predicate).toArray(JournalEntry[]::new);
  }
}