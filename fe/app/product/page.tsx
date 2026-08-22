import { DataTable } from "@/components/data-table";

import { columns } from "./column";
import { UploadSection } from "@/components/molecules";
import { getAllProducts, UploadProductExcel } from "../action/product";

export default async function Page() {
  return (
    <>
      <UploadSection uploadFn={UploadProductExcel} />
      <DataTable columns={columns} fetchFunction={getAllProducts}/>
    </>
  );
}
