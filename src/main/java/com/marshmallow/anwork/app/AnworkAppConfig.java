package com.marshmallow.anwork.app;

import java.io.File;
import java.util.function.Consumer;

/**
 * This is a data bucket for configuration options to be passed around through
 * the ANWORK app.
 *
 * <p>
 * Created Sep 11, 2017
 * </p>
 *
 * @author Andrew
 */
public class AnworkAppConfig {

  // These fields are set to their defaults.
  private String context = "default-context";
  private File persistenceRoot = new File(".");
  private boolean debug = false;
  private Consumer<String> debugPrinter = new Consumer<String>() {
    @Override
    public void accept(String string) {
      if (debug) {
        System.out.println("debug: " + string);
      }
    }
  };

  public String getContext() {
    return context;
  }

  public void setContext(String context) {
    this.context = context;
  }

  public File getPersistenceRoot() {
    return persistenceRoot;
  }

  public void setPersistenceRoot(File persistenceRoot) {
    this.persistenceRoot = persistenceRoot;
  }

  public boolean getDebug() {
    return debug;
  }

  public void setDebug(boolean debug) {
    this.debug = debug;
  }

  public Consumer<String> getDebugPrinter() {
    return debugPrinter;
  }

  public void setDebugPrinter(Consumer<String> debugPrinter) {
    this.debugPrinter = debugPrinter;
  }
}
