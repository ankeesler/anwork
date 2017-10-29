package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.Action;
import com.marshmallow.anwork.app.cli.ActionCreator;
import com.marshmallow.anwork.app.cli.ArgumentValues;
import com.marshmallow.anwork.journal.Journal;
import com.marshmallow.anwork.task.TaskManager;
import com.marshmallow.anwork.task.TaskManagerJournalEntry;

import java.util.Arrays;

/**
 * This is a {@link ActionCreator} for the journal commands in the ANWORK app.
 *
 * <p>
 * Created Oct 4, 2017
 * </p>
 *
 * @author Andrew
 */
public class JournalCliActionCreator implements ActionCreator {

  @Override
  public Action createAction(String commandName) {
    switch (commandName) {
      case "show-all":
        return new TaskManagerCliAction() {
          @Override
          public void run(AnworkAppConfig config,
                          ArgumentValues flags,
                          String[] args,
                          TaskManager manager) {
            System.out.println(Arrays.toString(manager.getJournal().getEntries()));
          }
        };
      case "show":
        return new TaskManagerCliAction() {
          @Override
          public void run(AnworkAppConfig config,
                          ArgumentValues flags,
                          String[] args,
                          TaskManager manager) {
            Journal<TaskManagerJournalEntry> journal = manager.getJournal(args[0]);
            if (journal == null) {
              System.out.println("No entries for task named " + args[0]);
            } else {
              System.out.println(Arrays.toString(journal.getEntries()));
            }
          }
        };
      default:
        return null; // error!
    }
  }
}
