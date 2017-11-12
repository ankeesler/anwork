package com.marshmallow.anwork.task.test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertTrue;

import com.marshmallow.anwork.journal.JournalEntry;
import com.marshmallow.anwork.task.Task;
import com.marshmallow.anwork.task.TaskManager;

/**
 * These are general utilities used when testing {@link TaskManager} instances.
 *
 * <p>
 * Created Nov 12, 2017
 * </p>
 *
 * @author Andrew
 */
public class TaskManagerTestUtilities {

  /**
   * This method asserts that the provided list of {@code expectedNames} matches what is actually
   * in the provided {@link TaskManager}'s {@link Journal}. The {@code expectedNames} should match
   * what is returned by {@link Task#getName}.
   *
   * @param expectedNames The names of {@link Task}'s as they appear in a {@link TaskManager}'s
   *     {@link Journal}
   */
  public static void assertJournalEntriesEqual(TaskManager manager, String...expectedNames) {
    // See TaskManagerJournalCache for more discussion on this weird logic.
    JournalEntry[] entries = manager.getJournal().getEntries();
    assertEquals(expectedNames.length, entries.length);
    for (int i = 0; i < expectedNames.length; i++) {
      assertTrue(entries[i].getTitle().contains(expectedNames[i]));
    }
  }

}
