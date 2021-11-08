import React, { useRef, useEffect } from "react";

const createRootElement = (id: string) => {
  const rootContainer = document.createElement("div");
  rootContainer.setAttribute("id", id);
  return rootContainer;
};

const addRootElement = (rootElem: Element) => {
  document.body.insertBefore(
    rootElem,
    document.body.lastElementChild?.nextElementSibling || null
  );
};

export const usePortal = (id: string) => {
  const rootElemRef = useRef<Element>(null);

  useEffect(() => {
    // Look for existing target dom element to append to
    const existingParent = document.querySelector(`#${id}`);
    // Parent is either a new root or the existing dom element
    const parentElem = existingParent || createRootElement(id);

    // If there is no existing DOM element, add a new one.
    if (!existingParent) {
      addRootElement(parentElem);
    }

    // Add the detached element to the parent
    // @ts-ignore
    parentElem.appendChild(rootElemRef.current);

    return () => {
      // @ts-ignore
      rootElemRef.current.remove();
      if (!parentElem.childElementCount) {
        parentElem.remove();
      }
    };
  }, [id]);

  const getRootElem = () => {
    if (!rootElemRef.current) {
      // @ts-ignore
      rootElemRef.current = document.createElement("div");
    }
    return rootElemRef.current;
  };

  return getRootElem();
};
