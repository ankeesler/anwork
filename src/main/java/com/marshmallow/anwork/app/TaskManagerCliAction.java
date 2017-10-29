package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.Action;
import com.marshmallow.anwork.app.cli.ArgumentValues;
import com.marshmallow.anwork.core.FilePersister;
import com.marshmallow.anwork.core.Persister;
import com.marshmallow.anwork.task.TaskManager;

import java.io.File;
import java.util.Collection;
import java.util.Collections;

/**
 * This is a {@link Action} that loads a {@link TaskManager}, does something, and
 * then saves the {@link TaskManager}.
 *
 * <p>
 * Created Sep 11, 2017
 * </p>
 *
 * @author Andrew
 */
public abstract class TaskManagerCliAction implements Action {

  @Override
  public void run(ArgumentValues flags, String[] args) {
    AnworkAppConfig config = new AnworkAppConfig(flags);
    try {
      TaskManager manager = loadTaskManager(config);
      run(config, flags, args, manager);
      saveTaskManager(config, manager);
    } catch (Exception e) {
      throw new IllegalStateException("Failed task manager action!", e);
    }
  }

  /**
   * Run the CLI action on a {@link TaskManager}.
   *
   * @param config The {@link AnworkAppConfig} applied to this {@link Action}
   * @param flags The {@link ArgumentValues} associated with this {@link Action}
   * @param args The CLI arguments
   * @param manager The task manager
   */
  public abstract void run(AnworkAppConfig config,
                           ArgumentValues flags,
                           String[] args,
                           TaskManager manager);

  private TaskManager loadTaskManager(AnworkAppConfig config) throws Exception {
    String context = config.getContext();
    File persistenceRoot = config.getPersistenceRoot();
    Persister<TaskManager> persister = new FilePersister<TaskManager>(persistenceRoot);
    if (!config.getDoPersist()) {
      config.getDebugPrinter().accept("not loading task manager because"
                                      + " persist command line option set to false");
      return new TaskManager();
    } else if (!persister.exists(context)) {
      config.getDebugPrinter().accept("context " + context
                                      + " does not exist at root " + persistenceRoot + "!"
                                      + " creating new task manager");
      return new TaskManager();
    }

    Collection<TaskManager> loadeds = persister.load(context, TaskManager.SERIALIZER);
    if (loadeds.size() != 1) {
      throw new IllegalStateException("Persistence root " + persistenceRoot
                                      + " and context " + context
                                      + " contains 0 or more than 1 task manager.");
    }
    return loadeds.toArray(new TaskManager[0])[0];
  }

  private void saveTaskManager(AnworkAppConfig config, TaskManager taskManager) throws Exception {
    Persister<TaskManager> persister = new FilePersister<TaskManager>(config.getPersistenceRoot());
    if (!config.getDoPersist()) {
      config.getDebugPrinter().accept("not saving task manager because"
                                      + " no-persist command line option set");
    } else {
      persister.save(config.getContext(),
                     TaskManager.SERIALIZER,
                     Collections.singleton(taskManager));
    }
  }
}
