package com.marshmallow.anwork.journal;

import java.util.ArrayList;
import java.util.List;

/**
 * This is a {@link Journal} implementation that stores its {@link JournalEntry}'s in a
 * {@link List}.
 *
 * <p>
 * Created Sep 17, 2017
 * </p>
 *
 * @author Andrew
 */
public class BaseJournal implements Journal {

  private final List<JournalEntry> entries = new ArrayList<JournalEntry>();

  @Override
  public void addEntry(JournalEntry entry) {
    entries.add(entry);
  }

  @Override
  public JournalEntry[] getEntries() {
    return entries.toArray(new JournalEntry[0]);
  }
}
