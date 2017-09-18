package com.marshmallow.anwork.task;

import com.marshmallow.anwork.journal.FilteredJournal;
import com.marshmallow.anwork.journal.Journal;

import java.util.HashMap;
import java.util.Map;

/**
 * This is a cache to store the {@link Journal} instances associated with {@link Task}'s managed by
 * the {@link TaskManager}.
 *
 * <p>
 * Created Sep 18, 2017
 * </p>
 *
 * @author Andrew
 */
public class TaskManagerJournalCache {

  private final Journal journal;
  // This is a map from task name to FilteredJournal.
  private final Map<String, Journal> cache = new HashMap<String, Journal>();

  /**
   * Create a cache given a root journal.
   *
   * @param journal The root journal to use in this cache
   */
  public TaskManagerJournalCache(Journal journal) {
    this.journal = journal;
  }

  /**
   * Get the {@link Journal} for a {@link Task} name, or create one if there is none that exists.
   *
   * @param taskName The name of the task
   */
  public Journal get(String taskName) {
    Journal journal = cache.get(taskName);
    if (journal == null) {
      // Note: this is pretty fuzzy logic. We should figure out a better filter for a per-task
      // journal.
      journal = new FilteredJournal(this.journal, (entry) -> entry.getTitle().contains(taskName));
      cache.put(taskName, journal);
    }
    return journal;
  }
}
