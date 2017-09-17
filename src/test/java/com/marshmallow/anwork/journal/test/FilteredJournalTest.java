package com.marshmallow.anwork.journal.test;

import static org.junit.Assert.assertEquals;

import com.marshmallow.anwork.journal.BaseJournal;
import com.marshmallow.anwork.journal.BaseJournalEntry;
import com.marshmallow.anwork.journal.FilteredJournal;
import com.marshmallow.anwork.journal.Journal;
import com.marshmallow.anwork.journal.JournalEntry;

import java.util.function.Predicate;

import org.junit.Before;
import org.junit.Test;

/**
 * This is a test for {@link FilteredJournal}.
 *
 * <p>
 * Created Sep 17, 2017
 * </p>
 *
 * @author Andrew
 */
public class FilteredJournalTest {

  private static void assertJournalEquals(Journal journal, JournalEntry[] expectedEntries) {
    JournalEntry[] actualEntries = journal.getEntries();
    assertEquals(expectedEntries.length, actualEntries.length);
    for (int i = 0; i < expectedEntries.length; i++) {
      assertEquals(expectedEntries[i].getTitle(), actualEntries[i].getTitle());
      assertEquals(expectedEntries[i].getDescription(), actualEntries[i].getDescription());
    }
  }

  private Journal journal = new BaseJournal();

  /**
   * Setup the journal under test with test entries.
   */
  @Before
  public void setupJournal() {
    journal.addEntry(new BaseJournalEntry("title:a", "description:a"));
    journal.addEntry(new BaseJournalEntry("title:b", "description:b"));
    journal.addEntry(new BaseJournalEntry("title:c", "description:c"));
  }

  @Test
  public void testEmpty() {
    assertEquals(0, new FilteredJournal(new BaseJournal(), entry -> true).getEntries().length);
  }

  @Test
  public void testAllMatchFilter() {
    assertJournalEquals(filter(entry -> entry.getTitle().startsWith("title")),
                        journal.getEntries());
  }

  @Test
  public void testNoneMatchFilter() {
    assertJournalEquals(filter(entry -> entry.getDescription().startsWith("title")),
                        new BaseJournalEntry[0]);
  }

  @Test
  public void testFilterPersistence() {
    Journal filteredJournal = filter(entry -> entry.getTitle().endsWith(":b"));
    assertJournalEquals(filteredJournal, new BaseJournalEntry[] {
        new BaseJournalEntry("title:b", "description:b"),
    });

    filteredJournal.addEntry(new BaseJournalEntry("title:d", "description:d"));
    filteredJournal.addEntry(new BaseJournalEntry("another title :b", "tuna fish marlin"));
    filteredJournal.addEntry(new BaseJournalEntry(":b_title", ":b_description"));
    assertJournalEquals(filteredJournal, new BaseJournalEntry[] {
        new BaseJournalEntry("title:b", "description:b"),
        new BaseJournalEntry("another title :b", "tuna fish marlin"),
    });

    // If we add entries to the BaseJournal, the events should still go through the filter.
    journal.addEntry(new BaseJournalEntry("title:e", "description:e"));
    journal.addEntry(new BaseJournalEntry(":b", "andrew"));
    assertJournalEquals(filteredJournal, new BaseJournalEntry[] {
        new BaseJournalEntry("title:b", "description:b"),
        new BaseJournalEntry("another title :b", "tuna fish marlin"),
        new BaseJournalEntry(":b", "andrew"),
    });    
  }

  private Journal filter(Predicate<JournalEntry> predicate) {
    return new FilteredJournal(journal, predicate);
  }
}
